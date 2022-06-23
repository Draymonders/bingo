create database scenes;

create table t_stock (
    id bigint primary key auto_increment comment '主键',
    product_id bigint default '0',
    stock_num bigint default '0'
) engine=InnoDB;

-- 创建product_id维度的key
alter table t_stock add index idx_product_id(`product_id`);