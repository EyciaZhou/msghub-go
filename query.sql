select * from msghub.msg order by SnapTime;
select count(*) from msghub.msg;

SELECT * FROM msghub.msg WHERE id="0000"  LIMIT  1;
SELECT
				(id, SnapTime, PubTime, SourceURL, Body, Title, SubTitle, CoverImg, ViewType, Frm, Tag, Topic)
			FROM msg
			WHERE id="001"
			LIMIT 1;