-- 重命名 config 表的字段
ALTER TABLE config RENAME COLUMN basePath TO base_path;
ALTER TABLE config RENAME COLUMN jwtSecret TO jwt_secret;
ALTER TABLE config RENAME COLUMN jwttokenKey TO jwt_token_key;
ALTER TABLE config RENAME COLUMN serverPort TO server_port;
ALTER TABLE config RENAME COLUMN maxRequestBodySize TO max_request_body_size;
ALTER TABLE config RENAME COLUMN createTime TO create_time;
ALTER TABLE config RENAME COLUMN updateTime TO update_time;
ALTER TABLE config RENAME COLUMN createUser TO create_user;
ALTER TABLE config RENAME COLUMN sortNo TO sortno;

-- 重命名 user 表的字段(user是SQL关键字,需要用引号括起来)
ALTER TABLE "user" RENAME COLUMN userName TO user_name;
ALTER TABLE "user" RENAME COLUMN chainType TO chain_type;
ALTER TABLE "user" RENAME COLUMN chainAddress TO chain_address;
ALTER TABLE "user" RENAME COLUMN createTime TO create_time;
ALTER TABLE "user" RENAME COLUMN updateTime TO update_time;
ALTER TABLE "user" RENAME COLUMN createUser TO create_user;
ALTER TABLE "user" RENAME COLUMN sortNo TO sortno;
ALTER TABLE "user" RENAME TO userinfo;

-- 重命名 category 表的字段
ALTER TABLE category RENAME COLUMN hrefURL TO href_url;
ALTER TABLE category RENAME COLUMN hrefTarget TO href_target;
ALTER TABLE category RENAME COLUMN templateFile TO template_file;
ALTER TABLE category RENAME COLUMN childTemplateFile TO child_template_file;
ALTER TABLE category RENAME COLUMN createTime TO create_time;
ALTER TABLE category RENAME COLUMN updateTime TO update_time;
ALTER TABLE category RENAME COLUMN createUser TO create_user;
ALTER TABLE category RENAME COLUMN sortNo TO sortno;

-- 重命名 content 表的字段
ALTER TABLE content RENAME COLUMN hrefURL TO href_url;
ALTER TABLE content RENAME COLUMN categoryID TO category_id;
ALTER TABLE content RENAME COLUMN categoryName TO category_name;
ALTER TABLE content RENAME COLUMN templateFile TO template_file;
ALTER TABLE content RENAME COLUMN signAddress TO sign_address;
ALTER TABLE content RENAME COLUMN signChain TO sign_chain;
ALTER TABLE content RENAME COLUMN txID TO tx_id;
ALTER TABLE content RENAME COLUMN createTime TO create_time;
ALTER TABLE content RENAME COLUMN updateTime TO update_time;
ALTER TABLE content RENAME COLUMN createUser TO create_user;
ALTER TABLE content RENAME COLUMN sortNo TO sortno;
-- 为 content 表添加 content_type 字段
ALTER TABLE content ADD COLUMN content_type INTEGER;
UPDATE content SET content_type = 0 WHERE markdown != '';
UPDATE content SET content_type = 1 WHERE markdown = '';

-- 重命名 site 表的字段
ALTER TABLE site RENAME COLUMN themePC TO theme_pc;
ALTER TABLE site RENAME COLUMN themeWAP TO theme_wap;
ALTER TABLE site RENAME COLUMN themeWX TO theme_wx;
ALTER TABLE site RENAME COLUMN createTime TO create_time;
ALTER TABLE site RENAME COLUMN updateTime TO update_time;
ALTER TABLE site RENAME COLUMN createUser TO create_user;
ALTER TABLE site RENAME COLUMN sortNo TO sortno;


-- 删除现有的触发器和虚拟表
DROP TRIGGER IF EXISTS trigger_content_insert;
DROP TRIGGER IF EXISTS trigger_content_delete;
DROP TRIGGER IF EXISTS trigger_content_update;
DROP TABLE IF EXISTS fts_content;


-- 使用设置-->更新SQL 的功能,按照步骤执行以下SQL,每次执行一行

-- 重新创建虚拟表(使用新的字段名)
CREATE VIRTUAL TABLE IF NOT EXISTS fts_content USING fts5(
		markdown, 
	    tokenize = 'simple 0',
		content='content', 
		content_rowid='rowid'
	);

-- 重新创建触发器(使用新的字段名)
CREATE TRIGGER IF NOT EXISTS trigger_content_insert AFTER INSERT ON content
BEGIN
    INSERT INTO fts_content (rowid, markdown) VALUES (new.rowid,  new.markdown);
END;

CREATE TRIGGER IF NOT EXISTS trigger_content_delete AFTER DELETE ON content
BEGIN
    INSERT INTO fts_content (fts_content, rowid, markdown) VALUES ('delete', old.rowid, old.markdown);
END;

CREATE TRIGGER IF NOT EXISTS trigger_content_update AFTER UPDATE ON content
BEGIN
    INSERT INTO fts_content (fts_content, rowid, markdown) VALUES ('delete', old.rowid,  old.markdown);
	INSERT INTO fts_content (rowid, markdown) VALUES (new.rowid, new.markdown);
END;


INSERT INTO fts_content (rowid, markdown) SELECT rowid, markdown FROM content;