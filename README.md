# Description

This service is taking metric **-m** from graphite server **-U**, authenticated by **-u** with token **-a** with timerange from **-range1From** untill **-range1Until** and compare it with timerange between **-range2From** and **-range2Until**  
If there is increasing more, than **-wi** it will exit with code 1 and if more, than **-ci** - with exit code 2.  
Same for decreasing with arguments **-wd** and **-cd**  
Also exit message will be percentage of the difference between time ranges.  

## Functions
This service will analyze metrics from graphite in one of next modes:
- percent (default):  
	takes mertic within range1 (from-until) and range2 (from-until)  
	and count as percent 2 metric from the 1st one (2nd*100/1st)
- absolute:  
	simular to percent, but with absolute values (2nd-1st)
- absoluteSingle:  
	tales metric only within range1 (from-until)
All of these methods print the result and give you exit code according nagios standards  
  
Warning limits include (intersect) threshold you specify (<= or >=) and criticals do not


# Usage
- -U string
        Base address of your graphite server e.g. https://graphite.protury.info/
- -u string
        User, which has rights to access Graphite
- -p string
        Password to access the graphite-API. For example 'qqq'
- -mode string
      	Mode of analysis of metrics. E.G. percentageDiff, absoluteDiff, absoluteCmp (default "percentageDiff")
- -ci float
      	Increasing. Metrics above this threshold will be marked as critical (default -0.10101)
- -wi float
      	Increasing. Metrics above this threshold will be marked as warning (default -0.10101)
- -cd float
    	Decreasing. Metrics below this threshold will be marked as critical (default -0.10101)
- -wd float
     	Decreasing. Metrics below this threshold will be marked as warning (default -0.10101)
- -m string
        Name of metric or metric filter e.g. qqqq.test.leoleovich.currentProblems
- -range1From int
    	Amount of seconds ago for the 1st range (from) (default 90000) = 1 day and 1 hour
- -range1Until int
    	Amount of seconds ago for the 1st range (until) (default 86400) = 1 day
- -range2From int
    	Amount of seconds ago for the 2st range (from) (default 3600) = 1 hour
- -range2Until int
    	Amount of seconds ago for the 2st range (until) (default 0) = current time
- -d	Debug mode will print a lot of additinal info

# Examples
- ```graphite2monitoring -U 'https://graphite.protury.info/' -mode 'absoluteCmp' -m 'sumSeries(offset(scale(all.about.money*,0),1))' -wd '-0.10101' -cd '1' -wi '3' -ci '-0.10101' -range1From '1800' -range1Until '0' -range2From '1800' -range2Until 0```  
Metric is ok (low limits (1.00 -0.10) < 2.00 < high limits (3.00 -0.10))  
```echo $?```  
0
- ```graphite2monitoring -U 'https://graphite.protury.info/' -mode 'absoluteCmp' -m 'sumSeries(offset(scale(all.about.money*,0),1))' -wd '-0.10101' -cd '1' -wi '2' -ci '-0.10101' -range1From '1800' -range1Until '0' -range2From '1800' -range2Until 0```  
Metric is above warning threshold (2.00 => 2.00)  
```echo $?```  
1
- ```graphite2monitoring -U 'https://graphite.protury.info/' -mode 'absoluteCmp' -m 'sumSeries(offset(scale(all.about.money*,0),1))' -wd '-0.10101' -cd '3' -wi '4' -ci '-0.10101' -range1From '1800' -range1Until '0' -range2From '1800' -range2Until 0```  
Metric is below critical threshold (2.00 < 3.00)  
```echo $?```  
2