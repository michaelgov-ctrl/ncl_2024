$log = Get-Content .\login.log

$hash = @{}
foreach ($l in $log) {
    $fields = $l -split '\s'
    $hash[$fields[3]] = @{}
}

$max = 0
foreach ($l in $log) {
    $fields = $l -split '\s'
    $hash[$fields[3]][$fields[2]]+=1
    
    if ($hash[$fields[3]].count -gt $max) {
        $max = $hash[$fields[3]].count
        $entryToCareAbout = [pscustomobject]@{
            user = $fields[3]
            ips = $hash[$fields[3]]
        }
    }
}

$entryToCareAbout