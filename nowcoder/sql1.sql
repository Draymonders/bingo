-- 计算各个视频的平均完播率
-- https://www.nowcoder.com/practice/96263162f69a48df9d84a93c71045753



-- 解法1： 分别求合法的播放数，总播放数，进行相除
select
    t1.video_id,
    format(IFNULL(completePlayCnt, 0) / allPlayCnt, 3) avg_comp_play_rate
from
(
    -- 查询所有的符合条件的视频
    select
        uv.video_id,
        count(uv.video_id) as allPlayCnt
    from tb_user_video_log as uv
    where uv.start_time >= '2021-01-01 00:00:00'
      and uv.end_time < '2022-01-01 00:00:00'
    group by uv.video_id
) as t1
left join
(
    -- 查询完成播放的视频
    select
        v.video_id video_id,
        count(uv.video_id) as completePlayCnt
    from tb_video_info as v
    inner join tb_user_video_log as uv
    on  uv.video_id = v.video_id
    where unix_timestamp(uv.end_time) - unix_timestamp(uv.start_time) >= v.duration
        and uv.start_time >= '2021-01-01 00:00:00'
        and uv.end_time < '2022-01-01 00:00:00'
    group by v.video_id
) as t2
on t1.video_id = t2.video_id
order by avg_comp_play_rate desc;


-- 解法2：直接使用sum求和
select
    v.video_id,
    round(sum(if(end_time - start_time >= duration, 1, 0)) / count(1), 3) as avg_comp_play_rate
from tb_video_info as v
inner join tb_user_video_log as uv
on v.video_id = uv.video_id and uv.start_time >= '2021-01-01 00:00:00'
group by v.video_id
order by avg_comp_play_rate desc;