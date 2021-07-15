-- Difference between prediction and actual temperature with difference of time between them
with hf as (
select
	datetime(
		created_at,
		'-'||strftime('%M', created_at)||' minutes',
		'-'||strftime('%S', created_at)||' seconds'
	) as created_at_utc,
	strftime('%Y-%m-%d %H:%M:%S', target_at) as target_at_utc,
	temperature,
	"condition",
	location,
	"source"
from hourly_forecasts
),
cf as (
	SELECT 
	datetime(
		created_at,
		'-'||strftime('%M', created_at)||' minutes',
		'-'||strftime('%S', created_at)||' seconds'
	) as created_at_utc,
	temperature,
	"condition",
	location,
	"source"
FROM current_forecasts
)
select
	hf."source" as "source",
	hf.location as location,
	hf.created_at_utc as prediction_date,
	hf.target_at_utc as target_date,
	hf.temperature as predicted_temperature,
	cf.temperature as actual_temperature,
	cf.temperature - hf.temperature as temperature_diff,
	hf."condition" as predicted_condition,
	cf."condition" as actual_condition,
	julianday(hf.target_at_utc) - julianday(hf.created_at_utc) as jd
from hf left join cf on
hf.target_at_utc = cf.created_at_utc AND 
hf.location = cf.location AND 
hf.source = cf.source
WHERE 
	prediction_date <> target_date AND 
	actual_temperature is not NULL;