create database lmbook;
create database lmbook_intr;
create database lmbook_article;
create database lmbook_user;
create database lmbook_payment;
create database lmbook_account;
create database lmbook_reward;
create database lmbook_comment;
create database lmbook_tag;

# 准备 canal 用户
CREATE USER 'canal'@'%' IDENTIFIED BY 'canal';
GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%' WITH GRANT OPTION;