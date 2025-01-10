-- 版本比较函数
DELIMITER //

CREATE FUNCTION compare_versions(version1 VARCHAR(255), version2 VARCHAR(255)) RETURNS INT
    READS SQL DATA
    DETERMINISTIC
BEGIN
    DECLARE v1_part INT;
    DECLARE v2_part INT;
    DECLARE v1_rest VARCHAR(255);
    DECLARE v2_rest VARCHAR(255);
    DECLARE pos1 INT;
    DECLARE pos2 INT;

    -- 初始化剩余字符串
    SET v1_rest = version1;
    SET v2_rest = version2;

    -- 循环比较每个部分
    WHILE LENGTH(v1_rest) > 0 OR LENGTH(v2_rest) > 0 DO
            -- 获取第一个部分
            SET pos1 = LOCATE('.', v1_rest);
            SET pos2 = LOCATE('.', v2_rest);

            -- 提取第一个部分的整数值
            SET v1_part = CAST(SUBSTRING(v1_rest, 1, IF(pos1 = 0, LENGTH(v1_rest), pos1 - 1)) AS UNSIGNED);
            SET v2_part = CAST(SUBSTRING(v2_rest, 1, IF(pos2 = 0, LENGTH(v2_rest), pos2 - 1)) AS UNSIGNED);

            -- 比较当前部分
            IF v1_part > v2_part THEN
                RETURN 1;
            ELSEIF v1_part < v2_part THEN
                RETURN -1;
            END IF;

            -- 更新剩余字符串
            SET v1_rest = IF(pos1 = 0, '', SUBSTRING(v1_rest, pos1 + 1));
            SET v2_rest = IF(pos2 = 0, '', SUBSTRING(v2_rest, pos2 + 1));
        END WHILE;

    -- 如果所有部分都相等，返回0
    RETURN 0;

END //

DELIMITER ;

SELECT compare_versions('12.2.3', '2.2.3'); -- 返回 1
SELECT compare_versions('1.2.3', '1.2.4'); -- 返回 -1
SELECT compare_versions('2.0.0', '1.9.9'); -- 返回 1
SELECT compare_versions('1.0.0', '1.0.0'); -- 返回 0