/*
 Navicat MySQL Data Transfer

 Source Server         : 127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 80013
 Source Host           : localhost:3306
 Source Schema         : reptiledata

 Target Server Type    : MySQL
 Target Server Version : 80013
 File Encoding         : 65001

 Date: 13/05/2025 21:44:11
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admin
-- ----------------------------
DROP TABLE IF EXISTS `admin`;
CREATE TABLE `admin`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `username` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `nickname` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '昵称',
  `password` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `salt` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '密码盐',
  `avatar` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '头像',
  `email` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '电子邮箱',
  `loginfailure` tinyint(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '失败次数',
  `logintime` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '登录时间',
  `createtime` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `updatetime` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
  `token` varchar(59) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT 'Session标识',
  `status` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'normal' COMMENT '状态',
  `canalId` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '所属渠道',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `username`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '管理员表' ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for admin_log
-- ----------------------------
DROP TABLE IF EXISTS `admin_log`;
CREATE TABLE `admin_log`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `admin_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
  `username` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '管理员名字',
  `url` varchar(1500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '操作页面',
  `title` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '日志标题',
  `content` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '内容',
  `ip` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT 'IP',
  `useragent` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT 'User-Agent',
  `createtime` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `name`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 415 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '管理员日志表' ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for attachment
-- ----------------------------
DROP TABLE IF EXISTS `attachment`;
CREATE TABLE `attachment`  (
  `id` int(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `admin_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
  `user_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '会员ID',
  `url` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '物理路径',
  `imagewidth` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '宽度',
  `imageheight` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '高度',
  `imagetype` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '图片类型',
  `imageframes` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '图片帧数',
  `filesize` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '文件大小',
  `mimetype` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT 'mime类型',
  `extparam` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '透传数据',
  `createtime` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建日期',
  `updatetime` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
  `uploadtime` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '上传时间',
  `storage` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'local' COMMENT '存储位置',
  `sha1` varchar(40) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '文件 sha1编码',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 22 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '附件表' ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for auth_group
-- ----------------------------
DROP TABLE IF EXISTS `auth_group`;
CREATE TABLE `auth_group`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `pid` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父组别',
  `name` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '组名',
  `rules` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '规则ID',
  `createtime` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `updatetime` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
  `status` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '状态',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '分组表' ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for auth_group_access
-- ----------------------------
DROP TABLE IF EXISTS `auth_group_access`;
CREATE TABLE `auth_group_access`  (
  `uid` int(10) UNSIGNED NOT NULL COMMENT '会员ID',
  `group_id` int(10) UNSIGNED NOT NULL COMMENT '级别ID',
  UNIQUE INDEX `uid_group_id`(`uid`, `group_id`) USING BTREE,
  INDEX `uid`(`uid`) USING BTREE,
  INDEX `group_id`(`group_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '权限分组表' ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for auth_rule
-- ----------------------------
DROP TABLE IF EXISTS `auth_rule`;
CREATE TABLE `auth_rule`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `type` enum('menu','file') CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'file' COMMENT 'menu为菜单,file为权限节点',
  `pid` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父ID',
  `name` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '规则名称',
  `title` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '规则名称',
  `icon` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '图标',
  `condition` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '条件',
  `remark` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `ismenu` tinyint(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否为菜单',
  `createtime` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `updatetime` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
  `weigh` int(10) NOT NULL DEFAULT 0 COMMENT '权重',
  `status` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '状态',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name`) USING BTREE,
  INDEX `pid`(`pid`) USING BTREE,
  INDEX `weigh`(`weigh`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 455 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '节点表' ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `pid` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父ID',
  `type` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '栏目类型',
  `name` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `nickname` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `flag` set('hot','index','recommend') CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `image` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '图片',
  `keywords` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '关键字',
  `description` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '描述',
  `diyname` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '自定义名称',
  `createtime` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
  `updatetime` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
  `weigh` int(10) NOT NULL DEFAULT 0 COMMENT '权重',
  `status` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '状态',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `weigh`(`weigh`, `id`) USING BTREE,
  INDEX `pid`(`pid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '分类表' ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for config
-- ----------------------------
DROP TABLE IF EXISTS `config`;
CREATE TABLE `config`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '变量名',
  `group` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '分组',
  `title` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '变量标题',
  `tip` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '变量描述',
  `type` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '类型:string,text,int,bool,array,datetime,date,file',
  `value` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '变量值',
  `content` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '变量字典数据',
  `rule` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '验证规则',
  `extend` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '扩展属性',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 26 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '系统配置' ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for domain
-- ----------------------------
DROP TABLE IF EXISTS `domain`;
CREATE TABLE `domain`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '网站名称',
  `domain` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '网站域名',
  `status` tinyint(2) NULL DEFAULT 1 COMMENT '状态（0：禁用，1：启用）',
  `create_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '添加时间',
  `expire_time` datetime(0) NULL DEFAULT NULL COMMENT '过期时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `domain`(`domain`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9376 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for domain_black
-- ----------------------------
DROP TABLE IF EXISTS `domain_black`;
CREATE TABLE `domain_black`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '网站名称',
  `domain` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `domain`(`domain`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 21860 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for domain_ranks
-- ----------------------------
DROP TABLE IF EXISTS `domain_ranks`;
CREATE TABLE `domain_ranks`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `serial_number` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '序列号',
  `domain_id` int(10) NULL DEFAULT 0 COMMENT '域名id',
  `keywords_id` int(10) NULL DEFAULT 0 COMMENT '关键字id',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '检索时间',
  `engines` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '',
  `type` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '',
  `ranks` int(5) NULL DEFAULT 0,
  `content` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `serial_number`(`serial_number`) USING BTREE,
  INDEX `domain_id`(`domain_id`) USING BTREE,
  INDEX `keywords_id`(`keywords_id`) USING BTREE,
  INDEX `create_time`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17273 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for keywords
-- ----------------------------
DROP TABLE IF EXISTS `keywords`;
CREATE TABLE `keywords`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `keywords` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '关键字',
  `level_id` int(10) NULL DEFAULT 0 COMMENT '等级id',
  `status` tinyint(2) NULL DEFAULT 1 COMMENT '状态（0：禁用，1：启用）',
  `create_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `search_time` datetime(0) NULL DEFAULT NULL COMMENT '搜索时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13750 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '搜索关键字' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for keywords_black
-- ----------------------------
DROP TABLE IF EXISTS `keywords_black`;
CREATE TABLE `keywords_black`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `keywords` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '关键字',
  `level_id` int(10) NULL DEFAULT 0 COMMENT '等级id',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for keywords_level
-- ----------------------------
DROP TABLE IF EXISTS `keywords_level`;
CREATE TABLE `keywords_level`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `level_name` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '关键字等级',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '关键字等级' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for keywords_status
-- ----------------------------
DROP TABLE IF EXISTS `keywords_status`;
CREATE TABLE `keywords_status`  (
  `keywords_id` int(10) NULL DEFAULT 0 COMMENT '关键字表id',
  `engines` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '搜索引擎名字',
  `type` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '搜索引擎类型',
  `status` tinyint(2) NULL DEFAULT 1 COMMENT '状态（0：禁用，1：等待检索，2：检索中，3：检索失败，4：已完成）',
  `query_id` varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '查询id',
  UNIQUE INDEX `kw_key`(`keywords_id`, `engines`, `type`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for search_engines
-- ----------------------------
DROP TABLE IF EXISTS `search_engines`;
CREATE TABLE `search_engines`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '搜索引擎名字',
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '搜索引擎英文名字',
  `url` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '搜索网址',
  `type` set('mobile','pc') CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'pc' COMMENT '引擎类型（1：pc，2：mobile）',
  `status` tinyint(2) NULL DEFAULT 1 COMMENT '状态（0：禁用，1启用）',
  `page` int(5) NULL DEFAULT 0 COMMENT '检索页数',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name`, `type`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '搜索引擎' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for search_realtime
-- ----------------------------
DROP TABLE IF EXISTS `search_realtime`;
CREATE TABLE `search_realtime`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `domain` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '域名',
  `keywords` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '关键字',
  `engines` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '搜索引擎',
  `type` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '类型',
  `page` int(5) NULL DEFAULT 1 COMMENT ' 检索页数',
  `status` tinyint(2) NULL DEFAULT 0 COMMENT '状态',
  `ranks` int(10) NULL DEFAULT 0 COMMENT '排名',
  `top_three` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT ' 前三',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for search_result
-- ----------------------------
DROP TABLE IF EXISTS `search_result`;
CREATE TABLE `search_result`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `serial_number` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '序列号',
  `engines` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '搜索引擎名字',
  `type` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '引擎类型',
  `keywords` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '关键字',
  `ranks` int(5) NULL DEFAULT 0 COMMENT '排名',
  `content` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '搜索结果',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT '查询时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `create_time`(`create_time`) USING BTREE,
  INDEX `keywords`(`keywords`) USING BTREE,
  INDEX `engines_id`(`type`, `engines`) USING BTREE,
  INDEX `ranks`(`ranks`) USING BTREE,
  INDEX `serial_number`(`serial_number`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '查询结果' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
