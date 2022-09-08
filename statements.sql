update metrics
set value = value + 1000
where id = 'testGauge';

select *
from metrics;