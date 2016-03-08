select * from msghub.msg order by SnapTime DESC;
select count(*) from msghub.msg;

select * from msghub.picref;

SELECT * FROM msghub.msg WHERE id="0000"  LIMIT  1;

SELECT
				Ref, Description, pid, nodenum
			FROM picref, pic_task_queue
			WHERE mid=1 AND picref.pid=pic_task_queue.id;

SELECT
				id, SnapTime, PubTime, SourceURL, Body, Title, SubTitle, CoverImg, ViewType, Frm, Tag, Topic
			FROM msg
			WHERE id="001"
			LIMIT 1;
            
SELECT
				id, SnapTime, PubTime, SourceURL, Title, SubTitle, CoverImg, ViewType, Frm, Tag, Topic
			FROM msg
			ORDER BY SnapTime DESC
			LIMIT 10;