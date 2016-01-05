This util is getting **metric** with timerange from **range1From** till **range1Until** and compare it with timerange from **range2From** till **range2Until**
If there is a difference (increasing/decreasing) more, than **-w** it will exit with code 1 and if more, than **-c** - with exit code 2.
Also exit message will be percentage of the difference between timeranges.