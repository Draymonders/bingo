package hive

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlGen_ByTableFull(t *testing.T) {
	a := assert.New(t)
	generator := &HiveGeneratorImpl{
		KeyPre:  "CF_1234",
		DB:      "db1",
		Table:   "tt1",
		Inputs:  []string{"k1", "k2"},
		Outputs: []string{"f1", "f2"},
	}

	{
		fmt.Println("==============hive to mq全量同步==============")
		toMqFull := generator.ToMqFull("cache")
		a.Equal(`SELECT
	'sink' AS _task,
    'cache' AS _dest,
	concat('CF_1234', '_', k1, '_', k2) AS _index,
	f1,
	f2
FROM
	db1.tt1
WHERE
	date = '${date}'
	`, toMqFull)
		fmt.Println(toMqFull)
	}

	{
		fmt.Println("==============hive to cache全量同步==============")
		toCacheFull := generator.ToCacheFull()
		//fmt.Println(tocacheFull)
		a.Equal(toCacheFull, `SELECT  
		concat('CF_1234', '_', T.k1, '_', T.k2) key,
        TO_JSON(
            named_struct(
				'f1', T.f1,
				'f2', T.f2
            )
        ) value
FROM    db1.tt1 T
WHERE   date = '${date}'
		`)
		fmt.Println(toCacheFull)
	}
}

func TestSqlGen_ByTableIncr(t *testing.T) {
	a := assert.New(t)
	generator := &HiveGeneratorImpl{
		KeyPre:  "CF_1234",
		DB:      "db1",
		Table:   "tt1",
		Inputs:  []string{"k1", "k2"},
		Outputs: []string{"f1", "f2"},
	}

	{
		fmt.Println("==============hive to mq增量同步==============")
		toMqIncr := generator.ToMqIncr("cache")
		a.Equal(`SELECT
    'sink' AS _task,
    'cache' AS _dest,
	concat('CF_1234', '_', T.k1, '_', T.k2) AS _index,
    T.f1 AS f1,
	T.f2 AS f2
FROM
    (
        SELECT
            k1,
			k2,
			f1,
			f2
        FROM
            db1.tt1
        where
            date = '${date}'
    ) T
    LEFT JOIN (
        SELECT
            k1,
			k2,
			f1,
			f2
        FROM
            db1.tt1
        where
            date = '${date-1}'
    ) L ON T.k1 = L.k1 AND T.k2 = L.k2
WHERE
    NOT T.f1 <=> L.f1
	OR NOT T.f2 <=> L.f2`, toMqIncr)
		fmt.Println(toMqIncr)
	}

	{
		fmt.Println("==============hive to cache增量同步==============")
		toCacheIncr := generator.ToCacheIncr()
		a.Equal(toCacheIncr, `SELECT  concat('CF_1234', '_', COALESCE(T.k1, L.k1), '_', COALESCE(T.k2, L.k2)) key,
        TO_JSON(
            named_struct(
				'f1', T.f1,
				'f2', T.f2
            )
        ) value
FROM    (
            SELECT  
					k1,
					k2,
					f1,
					f2
            FROM    db1.tt1
            WHERE   date = '${date}'
        ) T
FULL JOIN
        (
            SELECT  
					k1,
					k2,
					f1,
					f2
            FROM    db1.tt1
            WHERE   date = '${date-1}'
        ) L
ON		T.k1 = L.k1 AND T.k2 = L.k2
WHERE   
		NOT T.f1 <=> L.f1
		OR NOT T.f2 <=> L.f2
`)
		fmt.Println(toCacheIncr)
	}

}

func TestSqlGen_ByTableDecr(t *testing.T) {
	a := assert.New(t)
	generator := &HiveGeneratorImpl{
		KeyPre:  "CF_1234",
		DB:      "db1",
		Table:   "tt1",
		Inputs:  []string{"k1", "k2"},
		Outputs: []string{"f1", "f2"},
	}

	{
		//fmt.Println("==============hive to mq减量同步==============")

	}
	{
		fmt.Println("==============hive to cache减量同步==============")
		toCacheDecr := generator.ToCacheDecr()
		a.Equal(toCacheDecr, `SELECT  concat('CF_1234', '_', L.k1, '_', L.k2) key,
        '{}' value
FROM    (
            SELECT  
					k1,
					k2
            FROM    db1.tt1
            WHERE   date = '${date-1}'
        ) L
LEFT JOIN
        (
            SELECT  
					k1,
					k2
            FROM    db1.tt1
            WHERE   date = '${date}'
        ) T
ON     T.k1 = L.k1 AND T.k2 = L.k2
WHERE   
		T.k1 IS NULL
		AND T.k2 IS NULL`)
		fmt.Println(toCacheDecr)
	}

}

func TestSqlGen_ByTable(t *testing.T) {
	a := assert.New(t)
	generator := &HiveGeneratorImpl{
		KeyPre:  "CF_1234",
		DB:      "db1",
		Table:   "tt1",
		Inputs:  []string{"k1", "k2"},
		Outputs: []string{"f1", "f2"},
	}
	{
		fields := generator.genAllFields(",")
		a.Equal("k1,k2,f1,f2", fields)
	}
	{
		sql := generator.JudgePerm()
		a.Equal("select k1,k2,f1,f2 from db1.tt1 where date= '' limit 1", sql)
	}

	//fmt.Println(tocacheIncr)

}

func TestSqlGen_BySQL(t *testing.T) {
	a := assert.New(t)
	gen := &HiveGeneratorImpl{
		KeyPre:  "CF_1234",
		SQL:     "select company_id,user_id,product_id,order_id from ecom_governance.sentry_ds_redesign_1 where date = '${date}' and hour = '${hour}'",
		Inputs:  []string{"company_id", "user_id"},
		Outputs: []string{"product_id", "order_id"},
		Hourly:  true,
	}
	sql := gen.ToMqFull("cache")
	a.Equal(`SELECT
	'sink' AS _task,
    'cache' AS _dest,
	concat('CF_1234', '_', company_id, '_', user_id) AS _index,
	product_id,
	order_id
FROM
	(
		select company_id,user_id,product_id,order_id from ecom_governance.sentry_ds_redesign_1 where date = '${date}' and hour = '${hour}'
	)`, sql)

	sql = gen.ToMqIncr("cache")
	a.Equal(`SELECT
    'sink' AS _task,
    'cache' AS _dest,
	concat('CF_1234', '_', T.company_id, '_', T.user_id) AS _index,
    T.product_id AS product_id,
	T.order_id AS order_id
FROM
    (
        select company_id,user_id,product_id,order_id from ecom_governance.sentry_ds_redesign_1 where date = '${date}' and hour = '${hour}'
    ) T
    LEFT JOIN (
        select company_id,user_id,product_id,order_id from ecom_governance.sentry_ds_redesign_1 where date = '${date-1h}' and hour = '${hour-1}'
    ) L ON T.company_id = L.company_id AND T.user_id = L.user_id
WHERE
    NOT T.product_id <=> L.product_id
	OR NOT T.order_id <=> L.order_id`, sql)

	sql = gen.ToMqDecr("cache")
	a.Equal(`SELECT  
		'sink' AS _task,
		'cache' AS _dest,
		concat('CF_1234', '_', T.company_id, '_', T.user_id) AS _index,
		"" AS product_id,"" AS order_id
FROM    (
            select company_id,user_id,product_id,order_id from ecom_governance.sentry_ds_redesign_1 where date = '${date-1h}' and hour = '${hour-1}'
        ) L
LEFT JOIN
        (
            select company_id,user_id,product_id,order_id from ecom_governance.sentry_ds_redesign_1 where date = '${date}' and hour = '${hour}'
        ) T
ON     T.company_id = L.company_id AND T.user_id = L.user_id
WHERE   
		T.company_id IS NULL
		AND T.user_id IS NULL`, sql)

	sql = gen.ToCacheFull()
	a.Equal(`SELECT  
		concat('CF_1234', '_', T.company_id, '_', T.user_id) key,
        TO_JSON(
            named_struct(
				'product_id', T.product_id,
				'order_id', T.order_id
            )
        ) value
FROM    (
	select company_id,user_id,product_id,order_id from ecom_governance.sentry_ds_redesign_1 where date = '${date}' and hour = '${hour}'
) T`, sql)

	sql = gen.ToCacheIncr()
	a.Equal(`SELECT  concat('CF_1234', '_', COALESCE(T.company_id, L.company_id), '_', COALESCE(T.user_id, L.user_id)) key,
        TO_JSON(
            named_struct(
				'product_id', T.product_id,
				'order_id', T.order_id
            )
        ) value
FROM    (
            select company_id,user_id,product_id,order_id from ecom_governance.sentry_ds_redesign_1 where date = '${date}' and hour = '${hour}'
        ) T
FULL JOIN
        (
            select company_id,user_id,product_id,order_id from ecom_governance.sentry_ds_redesign_1 where date = '${date-1h}' and hour = '${hour-1}'
        ) L
ON		T.company_id = L.company_id AND T.user_id = L.user_id
WHERE   
		NOT T.product_id <=> L.product_id
		OR NOT T.order_id <=> L.order_id
`, sql)

	sql = gen.ToCacheDecr()
	a.Equal(`SELECT  concat('CF_1234', '_', L.company_id, '_', L.user_id) key,
        '{}' value
FROM    (
            select company_id,user_id,product_id,order_id from ecom_governance.sentry_ds_redesign_1 where date = '${date-1h}' and hour = '${hour-1}'
        ) L
LEFT JOIN
        (
            select company_id,user_id,product_id,order_id from ecom_governance.sentry_ds_redesign_1 where date = '${date}' and hour = '${hour}'
        ) T
ON     T.company_id = L.company_id AND T.user_id = L.user_id
WHERE   
		T.company_id IS NULL
		AND T.user_id IS NULL`, sql)
}

func TestDateReg(t *testing.T) {
	a := assert.New(t)
	sql := `... date = ${date} or date=${ date } or date = ${date-1} and or = ${date +1} or date = '${date-1h}'`
	submatch := dateReg.FindAllStringSubmatch(sql, -1)
	a.Equal(5, len(submatch))
}

func Test_SqlGen_BySQLFull(t *testing.T) {
	// a := assert.New(t)
	gen := &HiveGeneratorImpl{
		KeyPre: "creator_attr_jb_ccr",
		SQL: `
select author_id, score, "1" as attr_value
from dm_temai.live_streaming_ccr_author_score_v1
where date = '${date}'
and author_level in ('L0','L1','L2','L3','L4')
and score > 0.4
`,
		Inputs:  []string{"author_id"},
		Outputs: []string{"author_id", "score", "attr_value"},
		Hourly:  false,
	}

	fullSqlStr := gen.ToCacheFull()
	fmt.Println(fullSqlStr)
	fmt.Println("====")

	incrSqlStr := gen.ToCacheIncr()
	fmt.Println(incrSqlStr)

	fmt.Println("====")

	decrSqlStr := gen.ToCacheDecr()
	fmt.Println(decrSqlStr)
}
