# Description

This service is taking metric **-m** from graphite server **-U**, authenticated by **-u** with token **-a** with timerange from **-range1From** untill **-range1Until** and compare it with timerange between **-range2From** and **-range2Until**
If there is increasing more, than **-wi** it will exit with code 1 and if more, than **-ci** - with exit code 2.
Same for decreasing with arguments **-wd** and **-cd**
Also exit message will be percentage of the difference between timeranges.

# Usage
- -U string
        Base address of your graphite server e.g. https://graphite.protury.info/
- -a string
        AuthToken to access the graphite-API. For example 'qqq'  
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
- graphite2monitoring -U 'https://graphite.protury.info/' -a 'verySecretKey' -m 'all.about.money' -wd '0' -cd '0' -wi '5' -ci '10' -range1From '3600' -range1Until '3000' -range2From '600' -range2Until 0  
Decreasing of metric is: 78.7%  
echo $?  
0
- graphite2monitoring -U 'https://graphite.protury.info/' -a 'verySecretKey' -m 'all.about.money' -wd '70' -cd '80' -wi '5' -ci '10' -range1From '3600' -range1Until '3000' -range2From '600' -range2Until 0  
Decreasing of metric is: 78.7%  
echo $?  
1
- graphite2monitoring -U 'https://graphite.protury.info/' -a 'verySecretKey' -m 'all.about.money' -wd '5' -cd '10'
Increasing of metric is: 78.9%  
echo $?  
2
