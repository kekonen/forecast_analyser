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
FROM current_forecasts;


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
from hourly_forecasts;

select
	datetime(
		created_at,
		'-'||strftime('%M', created_at)||' minutes',
		'-'||strftime('%S', created_at)||' seconds'
	) as created_at_utc,
	datetime(
		target_at,
		SUBSTRING(target_at, 20, 3) || ' hours',
		'start of day'
	)
	as target_at_utc,
	temperature_min,
	temperature_max,
	"condition",
	location,
	"source"
from daily_forecasts;
