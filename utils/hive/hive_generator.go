package hive

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	dateReg = regexp.MustCompile("\\$\\{\\s*(date)[0-9a-zA-Z+\\-\\s]*}")
	hourReg = regexp.MustCompile("\\$\\{\\s*(hour)[0-9a-zA-Z+\\-\\s]*}")
)

type HiveGenerator interface {
	ToMqFull(dest string) string // mq 全量
	ToMqIncr(dest string) string // mq 增量
	ToMqDecr(dest string) string // mq 减量
	ToCacheFull() string         // hive -> Cache 全量
	ToCacheDecr() string         // hive -> Cache 减量
	ToCacheIncr() string         // hive -> Cache 增量

	ToCache() string // hive -> Cache任务

	JudgePerm() string
}

type HiveGeneratorImpl struct {
	KeyPre  string // CF_123
	DB      string
	Table   string
	SQL     string // 自定义SQL模式
	Inputs  []string
	Outputs []string
	Hourly  bool   // 是否小时级别更新
	Version string // 模型版本

	MQModelInfo    *MQModelInfo    // 模型信息，MQ同步使用
	CacheModelInfo *CacheModelInfo // 缓存信息，缓存同步使用
}

func (g *HiveGeneratorImpl) ToCache() string {
	return fmt.Sprintf(`
select key, value
from 
(
	%s
) 
union all (
	%s
)
`, g.ToCacheFull(), g.ToCacheDecr())
}

type MQModelInfo struct {
	ModelEnName string
	ModelCnName string
	Version     string
	Score       *float64
}

type CacheModelInfo struct {
	Score *float64
}

func (g *HiveGeneratorImpl) ToMqFull(dest string) string {
	if g.SQL != "" {
		return g.hive2MqFullBySQL(dest)
	}
	return g.hive2MqFullByTable(dest)
}

//注意：外部触发调度，${date}取值是T；对于执行频率为天/周/月的任务，${date}取值为T-1。参考：https://site.bytedance.net/docs/2107/2638/110632/
//hive2mq有两个任务：全量任务和增量任务，全量任务仅执行一次，增量任务每天都执行
//hive2Cache有两种模式：全量模式和增量模式
//	全量模式有两个任务：全量任务和减量任务，两个任务每天都会执行。减量任务用于清理T分区没有的数据
//	增量模式有两个任务：全量任务和增量任务，与hive2mq类似

/*
* hive2mq全量同步任务
SELECT

	    'sink' AS _task,
	    'Cache' AS _dest,
		concat('CF_1234', '_', A.shop_id) AS _index,
	    toutiao_id AS shop_toutiao_id,
	    init_time AS shop_init_time

FROM

	ecom.app_govern_sentry_factor_shop_margin_recharge_withdrawal_df

WHERE

	date = '${date}'
*/
func (g *HiveGeneratorImpl) hive2MqFullByTable(dest string) string {
	isHour := ""
	if g.Hourly {
		isHour = "AND hour = '${hour}'"
	}
	hiveSql := fmt.Sprintf(`SELECT
	'sink' AS _task,
    '%s' AS _dest,
	%s AS _index,
	%s
FROM
	%s.%s
WHERE
	date = '${date}'
	%s`, dest, g.genIndex(""), g.genFields(g.Outputs, ",\n\t"), g.DB, g.Table, isHour)
	return hiveSql
}

func (g *HiveGeneratorImpl) ToMqIncr(dest string) string {
	if g.SQL != "" {
		return g.hive2MqIncrBySQL(dest)
	}
	return g.hive2MqIncrByTable(dest)
}

/*
* hive2mq增量同步任务
SELECT

	    'sink' AS _task,
	    'Cache/redis' AS _dest,
		concat('CF_1234', '_', A.shop_id) AS _index,
	    T.toutiao_id AS toutiao_id,
	    T.init_time AS init_time

FROM

	(
	    SELECT
	        shop_id,
	        toutiao_id,
	        init_time
	    FROM
	        ecom.app_govern_sentry_factor_shop_margin_recharge_withdrawal_df
	    where
	        date = '${date}'
	) T
	LEFT JOIN (
	    SELECT
	        shop_id,
	        toutiao_id,
	        init_time
	    FROM
	        ecom.app_govern_sentry_factor_shop_margin_recharge_withdrawal_df
	    where
	        date = '${date-1}'
	) L ON T.shop_id = L.shop_id

WHERE

	NOT T.toutiao_id <=> L.toutiao_id
	OR NOT T.init_time <=> L.init_time
*/
func (g *HiveGeneratorImpl) hive2MqIncrByTable(dest string) string {
	dbTable := g.DB + "." + g.Table
	allFields := g.genAllFields(",\n\t\t\t")
	onCond := g.genJoinCondition()
	notEqualCond := g.notEqualCondition("\n\tOR ")
	if notEqualCond == "" {
		notEqualCond = "TRUE"
	}
	this := `date = '${date}'`   // 当前时刻
	last := `date = '${date-1}'` // 上一时刻
	if g.Hourly {
		this = `date = '${date}' AND hour = '${hour}'`
		last = `date = '${date-1h}' AND hour = '${hour-1}'`
	}
	hiveSql := fmt.Sprintf(`SELECT
    'sink' AS _task,
    '%s' AS _dest,
	%s AS _index,
    %s
FROM
    (
        SELECT
            %s
        FROM
            %s
        where
            %s
    ) T
    LEFT JOIN (
        SELECT
            %s
        FROM
            %s
        where
            %s
    ) L ON %s
WHERE
    %s`, dest, g.genIndex("T."), g.genAs(g.Outputs, ",\n\t"), allFields, dbTable, this, allFields, dbTable, last, onCond, notEqualCond)
	return hiveSql
}

// hive2mq减量任务 将所有字段置为空串
func (g *HiveGeneratorImpl) ToMqDecr(dest string) string {
	if g.SQL != "" {
		return g.hive2MqDecrBySQL(dest)
	}
	return g.hive2MqDecrByTable(dest)
}
func (g *HiveGeneratorImpl) hive2MqDecrByTable(dest string) string {
	db, table := g.DB, g.Table
	dbTable := db + "." + table
	fields := g.genFields(g.Inputs, ",\n\t\t\t\t\t")
	this := `date = '${date}'`   // 当前时刻
	last := `date = '${date-1}'` // 上一时刻
	if g.Hourly {
		this = `date = '${date}' AND hour = '${hour}'`
		last = `date = '${date-1h}' AND hour = '${hour-1}'`
	}
	sql := fmt.Sprintf(`SELECT  
		'sink' AS _task,
		'%s' AS _dest,
		%s AS _index,
		%s
FROM    (
            SELECT  
					%s
            FROM    %s
            WHERE   %s
        ) L
LEFT JOIN
        (
            SELECT  
					%s
            FROM    %s
            WHERE   %s
        ) T
ON     %s
WHERE   
		%s`,
		dest, g.genIndex("T."), g.genEmptyAs(g.Outputs), fields, dbTable, last, fields, dbTable, this, g.genJoinCondition(),
		g.isNullCondition(g.Inputs, "T.", "\n\t\tAND "))
	return sql
}

func (g *HiveGeneratorImpl) genEmptyAs(params []string) string {
	builder := strings.Builder{}
	for i, p := range params {
		builder.WriteString(fmt.Sprintf(`"" AS %s`, p))
		if i != len(params)-1 {
			builder.WriteString(",")
		}
	}
	return builder.String()
}

// T.field as field,
func (g *HiveGeneratorImpl) genAs(params []string, separator string) string {
	builder := strings.Builder{}
	for i, p := range params {
		builder.WriteString(fmt.Sprintf(`T.%s AS %s`, p, p))
		if i != len(params)-1 {
			builder.WriteString(separator)
		}
	}
	return builder.String()
}

// concat('CF_1234', '_', T.app_id, '_', T.user_id)
func (g *HiveGeneratorImpl) genIndex(prefix string) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(`concat('%s'`, g.KeyPre))
	for _, p := range g.Inputs {
		builder.WriteString(fmt.Sprintf(`, '_', %s%s`, prefix, p))
	}
	builder.WriteString(`)`)
	return builder.String()
}

func (g *HiveGeneratorImpl) genFields(params []string, separator string) string {
	builder := strings.Builder{}
	for i, p := range params {
		builder.WriteString(p)
		if i != len(params)-1 {
			builder.WriteString(separator)
		}
	}
	if g.MQModelInfo != nil {
		builder.WriteString(separator)

		builder.WriteString(fmt.Sprintf("'%v' as model_en_name", g.MQModelInfo.ModelEnName))
		builder.WriteString(separator)

		builder.WriteString(fmt.Sprintf("'%v' as model_cn_name", g.MQModelInfo.ModelCnName))
		builder.WriteString(separator)

		builder.WriteString(fmt.Sprintf("'%v' as version", g.MQModelInfo.Version))

		if g.MQModelInfo.Score != nil {
			builder.WriteString(separator)
			builder.WriteString(fmt.Sprintf("%v as score", *g.MQModelInfo.Score))
		}
	}
	return builder.String()
}

func (g *HiveGeneratorImpl) genAllFields(separator string) string {
	// 去重
	allFields := make([]string, 0, len(g.Inputs))
	added := make(map[string]bool)
	for i, p := range g.Inputs {
		added[p] = true
		allFields = append(allFields, g.Inputs[i])
	}
	for i, f := range g.Outputs {
		if !added[f] {
			allFields = append(allFields, g.Outputs[i])
		}
	}
	return g.genFields(allFields, separator)
}

// T.key1 = L.key1 AND T.key2 = L.key2
func (g *HiveGeneratorImpl) genJoinCondition() string {
	inParams := g.Inputs
	builder := strings.Builder{}
	for i, p := range inParams {
		builder.WriteString(fmt.Sprintf(`T.%s = L.%s`, p, p))
		if i != len(inParams)-1 {
			builder.WriteString(" AND ")
		}
	}
	return builder.String()
}

// NOT T.field1 <=> L.field1
func (g *HiveGeneratorImpl) notEqualCondition(separator string) string {
	params := g.Outputs
	builder := strings.Builder{}
	for i, p := range params {
		builder.WriteString(fmt.Sprintf(`NOT T.%s <=> L.%s`, p, p))
		if i != len(params)-1 {
			builder.WriteString(separator)
		}
	}
	return builder.String()
}

// L.key1 IS NULL
func (g *HiveGeneratorImpl) isNullCondition(params []string, prefix string, separator string) string {
	builder := strings.Builder{}
	for i, p := range params {
		builder.WriteString(fmt.Sprintf(`%s%s IS NULL`, prefix, p))
		if i != len(params)-1 {
			builder.WriteString(separator)
		}
	}
	return builder.String()
}

func (g *HiveGeneratorImpl) ToCacheFull() string {
	if g.SQL != "" {
		return g.hive2CacheFullBySQL()
	}
	return g.hive2CacheFullByTable()
}

/*
* hive2Cache全量同步任务
SELECT  concat('CF_**', '_', shop_id) key,

	TO_JSON(
	    named_struct(
	        'init_time', T.init_time
	    )
	) value

FROM    ecom.app_govern_sentry_factor_shop_margin_recharge_withdrawal_df
WHERE   date = '${date}'
*/
func (g *HiveGeneratorImpl) hive2CacheFullByTable() string {
	db, table := g.DB, g.Table
	dbTable := db + "." + table
	isHour := ""
	if g.Hourly {
		isHour = "AND hour = '${hour}"
	}
	sql := fmt.Sprintf(`SELECT  
		%s key,
        %s value
FROM    %s T
WHERE   date = '${date}'
		%s`, g.CacheKey("T."), g.CacheValue(), dbTable, isHour)
	return sql
}

func (g *HiveGeneratorImpl) ToCacheDecr() string {
	if g.SQL != "" {
		return g.hive2CacheDecrBySQL()
	}
	return g.hive2CacheDecrByTable()
}

/*
* hive2Cache减量任务，配合全量同步任务：找到T-1分区有，T分区没有的记录
SELECT  concat('CF_**', '_', L.shop_id) key,

	'{}' value

FROM    (

	    SELECT  shop_id
	    FROM    ecom.app_govern_sentry_factor_shop_margin_recharge_withdrawal_df
	    WHERE   date = '${date-1}'
	) L

LEFT JOIN

	(
	    SELECT  shop_id
	    FROM    ecom.app_govern_sentry_factor_shop_margin_recharge_withdrawal_df
	    WHERE   date = '${date}'
	) T

ON      L.shop_id = T.shop_id
WHERE   T.shop_id IS NULL
*/
func (g *HiveGeneratorImpl) hive2CacheDecrByTable() string {
	db, table := g.DB, g.Table
	dbTable := db + "." + table
	fields := g.genFields(g.Inputs, ",\n\t\t\t\t\t")
	this := `date = '${date}'`   // 当前时刻
	last := `date = '${date-1}'` // 上一时刻
	if g.Hourly {
		this = `date = '${date}' AND hour = '${hour}'`
		last = `date = '${date-1h}' AND hour = '${hour-1}'`
	}
	sql := fmt.Sprintf(`SELECT  %s key,
        '{}' value
FROM    (
            SELECT  
					%s
            FROM    %s
            WHERE   %s
        ) L
LEFT JOIN
        (
            SELECT  
					%s
            FROM    %s
            WHERE   %s
        ) T
ON     %s
WHERE   
		%s`,
		g.CacheKey("L."), fields, dbTable, last, fields, dbTable, this, g.genJoinCondition(), g.isNullCondition(g.Inputs, "T.", "\n\t\tAND "))
	return sql
}

func (g *HiveGeneratorImpl) ToCacheIncr() string {
	if g.SQL != "" {
		return g.hive2CacheIncrBySQL()
	}
	return g.hive2CacheIncrByTable()
}

/*
* hive2Cache增量同步任务
SELECT  concat('CF_**', '_', COALESCE(T.shop_id, L.shop_id)) key,

	TO_JSON(
	    named_struct(
	        'init_time', T.init_time
	    )
	) value

FROM    (

	            SELECT
						shop_id,
	                    init_time
	            FROM    ecom.app_govern_sentry_factor_shop_margin_recharge_withdrawal_df
	            WHERE   date = '${date}'
	        ) T

FULL JOIN

	        (
	            SELECT
						shop_id,
	                    init_time
	            FROM    ecom.app_govern_sentry_factor_shop_margin_recharge_withdrawal_df
	            WHERE   date = '${date-1}'
	        ) L

ON      T.shop_id = L.shop_id
WHERE

	NOT T.init_time <=> L.init_time
*/
func (g *HiveGeneratorImpl) hive2CacheIncrByTable() string {
	db, table := g.DB, g.Table
	dbTable := db + "." + table
	fields := g.genAllFields(",\n\t\t\t\t\t")
	joinCond := g.genJoinCondition()
	whereCond := g.notEqualCondition("\n\t\tOR ")
	this := `date = '${date}'`   // 当前时刻
	last := `date = '${date-1}'` // 上一时刻
	if g.Hourly {
		this = `date = '${date}' AND hour = '${hour}'`
		last = `date = '${date-1h}' AND hour = '${hour-1}'`
	}
	sql := fmt.Sprintf(`SELECT  %s key,
        %s value
FROM    (
            SELECT  
					%s
            FROM    %s
            WHERE   %s
        ) T
FULL JOIN
        (
            SELECT  
					%s
            FROM    %s
            WHERE   %s
        ) L
ON		%s
WHERE   
		%s
`,
		g.CacheKeyIncr(), g.CacheValue(), fields, dbTable, this, fields, dbTable, last, joinCond, whereCond)
	return sql
}

// concat('CF_', id, '_', T.shop_id)
func (g *HiveGeneratorImpl) CacheKey(prefix string) string {
	pks := g.Inputs
	builder := strings.Builder{}
	builder.WriteString(`concat('`)
	builder.WriteString(g.KeyPre)
	builder.WriteString(`'`)
	for _, pk := range pks {
		builder.WriteString(`, '_', `)
		builder.WriteString(fmt.Sprintf(`%s%s`, prefix, pk))
	}
	if len(g.Version) > 0 {
		builder.WriteString(`, '_', `)
		builder.WriteString(g.Version)
	}
	builder.WriteString(`)`)
	return builder.String()
}

// concat('CF_', id, '_', COALESCE(T.shop_id, L.shop_id))
func (g *HiveGeneratorImpl) CacheKeyIncr() string {
	pks := g.Inputs
	builder := strings.Builder{}
	builder.WriteString(`concat('`)
	builder.WriteString(g.KeyPre)
	builder.WriteString(`'`)
	for _, pk := range pks {
		builder.WriteString(`, '_', `)
		builder.WriteString(fmt.Sprintf(`COALESCE(T.%s, L.%s)`, pk, pk))
	}
	if len(g.Version) > 0 {
		builder.WriteString(`, '_', `)
		builder.WriteString(g.Version)
	}
	builder.WriteString(`)`)
	return builder.String()
}

func (g *HiveGeneratorImpl) CacheValue() string {
	fields := g.Outputs
	builder := strings.Builder{}
	for i, f := range fields {
		builder.WriteString(fmt.Sprintf("'%s', T.%s", f, f))
		if i != len(fields)-1 {
			builder.WriteString(",\n\t\t\t\t")
		}
	}
	if g.CacheModelInfo != nil {
		builder.WriteString(",\n\t\t\t\t")

		builder.WriteString(fmt.Sprintf("'%s', 1", "attr_value"))
		builder.WriteString(",\n\t\t\t\t")

		if g.CacheModelInfo.Score != nil {
			builder.WriteString(fmt.Sprintf("'score', T.%v", *g.CacheModelInfo.Score))
			builder.WriteString(",\n\t\t\t\t")
		}

		builder.WriteString("'date', '${date}'")
	}
	return fmt.Sprintf(`TO_JSON(
            named_struct(
				%s
            )
        )`, builder.String())
}

func (g *HiveGeneratorImpl) hive2MqFullBySQL(dest string) string {
	hiveSql := fmt.Sprintf(
		`SELECT
	%s AS _index,
	%s
FROM
	(
		%s
	)`, g.genIndex(""), g.genFields(g.Outputs, ",\n\t"), g.SQL)
	return hiveSql
}

func (g *HiveGeneratorImpl) hive2MqIncrBySQL(dest string) string {
	onCond := g.genJoinCondition()
	notEqualCond := g.notEqualCondition("\n\tOR ")
	if notEqualCond == "" {
		notEqualCond = "TRUE"
	}
	hiveSql := fmt.Sprintf(`SELECT
    'sink' AS _task,
    '%s' AS _dest,
	%s AS _index,
    %s
FROM
    (
        %s
    ) T
    LEFT JOIN (
        %s
    ) L ON %s
WHERE
    %s`, dest, g.genIndex("T."), g.genAs(g.Outputs, ",\n\t"), g.SQL, g.lastPeriod(g.SQL), onCond, notEqualCond)
	return hiveSql
}

func (g *HiveGeneratorImpl) hive2MqDecrBySQL(dest string) string {
	sql := fmt.Sprintf(`SELECT  
		'sink' AS _task,
		'%s' AS _dest,
		%s AS _index,
		%s
FROM    (
            %s
        ) L
LEFT JOIN
        (
            %s
        ) T
ON     %s
WHERE   
		%s`,
		dest, g.genIndex("T."), g.genEmptyAs(g.Outputs), g.lastPeriod(g.SQL), g.SQL, g.genJoinCondition(),
		g.isNullCondition(g.Inputs, "T.", "\n\t\tAND "))
	return sql
}

func (g *HiveGeneratorImpl) hive2CacheFullBySQL() string {
	sql := fmt.Sprintf(`SELECT  
		%s key,
        %s value
FROM    (
	%s
) T`, g.CacheKey("T."), g.CacheValue(), g.SQL)
	return sql
}

func (g *HiveGeneratorImpl) hive2CacheDecrBySQL() string {
	sql := fmt.Sprintf(`SELECT  %s key,
        '{}' value
FROM    (
            %s
        ) L
LEFT JOIN
        (
            %s
        ) T
ON     %s
WHERE   
		%s`,
		g.CacheKey("L."), g.lastPeriod(g.SQL), g.SQL, g.genJoinCondition(), g.isNullCondition(g.Inputs, "T.", "\n\t\tAND "))
	return sql
}

func (g *HiveGeneratorImpl) hive2CacheIncrBySQL() string {
	joinCond := g.genJoinCondition()
	whereCond := g.notEqualCondition("\n\t\tOR ")
	sql := fmt.Sprintf(`SELECT  %s key,
        %s value
FROM    (
            %s
        ) T
FULL JOIN
        (
            %s
        ) L
ON		%s
WHERE   
		%s
`,
		g.CacheKeyIncr(), g.CacheValue(), g.SQL, g.lastPeriod(g.SQL), joinCond, whereCond)
	return sql
}

// 获取上一个周期的同步sql，天级别：前一天，小时级别：前一个小时
func (g *HiveGeneratorImpl) lastPeriod(sql string) string {
	if g.Hourly {
		sql = dateReg.ReplaceAllStringFunc(sql, func(s string) string {
			return strings.ReplaceAll(s, "date", "date-1h")
		})
		return hourReg.ReplaceAllStringFunc(sql, func(s string) string {
			return strings.ReplaceAll(s, "hour", "hour-1")
		})
	} else {
		return dateReg.ReplaceAllStringFunc(sql, func(s string) string {
			return strings.ReplaceAll(s, "date", "date-1")
		})
	}
	//return strings.ReplaceAll(sql, "${date}", "${date-1}")
}

func (g *HiveGeneratorImpl) JudgePerm() string {
	dbName, tableName := g.DB, g.Table
	fields := g.genAllFields(",")
	return fmt.Sprintf(`select %s from %s.%s where date= '' limit 1`, fields, dbName, tableName)
}
