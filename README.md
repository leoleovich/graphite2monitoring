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


# Usage
- -U string
        Base address of your graphite server e.g. https://graphite.protury.info/
- -a string
        AuthToken to access the graphite-API. For example 'qqq'
- -mode string
      	Mode of analysis of metrics. E.G. percentageDiff, absoluteDiff, absoluteCmp (default "percentageDiff")
- -ci int
        Metrics above this threshold will be marked as critical (default 40)
- -wi int
        Metrics above this threshold will be marked as warning (default 20)
- -cd int
        Metrics below this threshold will be marked as critical (default 40)
- -wd int
        Metrics below this threshold will be marked as warning (default 20)
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
- -u string
    	User, which has rights to access Graphite (default "graphite")
- -d	Debug mode will print a lot of additinal info

# Examples
- ```graphite2monitoring -U 'https://graphite.protury.info/' -a 'verySecretKey' -m 'all.about.money' -wd '0' -cd '0' -wi '5' -ci '10' -range1From '3600' -range1Until '3000' -range2From '600' -range2Until 0```  
Decreasing of metric is: 78.7%  
```echo $?```  
0
- ```graphite2monitoring -U 'https://graphite.protury.info/' -a 'verySecretKey' -m 'all.about.money' -wd '70' -cd '80' -wi '5' -ci '10' -range1From '3600' -range1Until '3000' -range2From '600' -range2Until 0```  
Decreasing of metric is: 78.7%  
```echo $?```  
1
- ```graphite2monitoring -U 'https://graphite.protury.info/' -a 'verySecretKey' -m 'all.about.money' -wd '5' -cd '10'  
Increasing of metric is: 78.9%  
```echo $?```  
2
- ```graphite2monitoring  -U 'https://graphite.protury.info/' -a 'verySecretKey' -m 'all.about.money' -wd '2' -cd '1' -wi '5' -ci '10' -range1From '3600' -range1Until '3000' -range2From '600' -range2Until 0 -mode 'singleAbsolute'```  
Metric is below critical threshold (0.000000 < 1.000000)  
```echo $?```  
2
- ```graphite2monitoring  -U 'https://graphite.innogames.de/' -a 'wtL$cDaCB%PpqtofjPNj1Ux+' -m 'nagios.wallOfShame.oleg_obleukhov.currentProblems' -wd '0' -cd '-1' -wi '5' -ci '10' -range1From '3600' -range1Until '3000' -range2From '600' -range2Until 0 -mode 'singleAbsolute'```  
Metric is ok (low limits (-1.000000 0.000000) < 0.000000 < high limits (5.000000 10.000000))  
```echo $?```
0