
function gop ($_, $in){
    switch ($_) {
        "qr" {
            .\Goshh-Client.exe -qr
        }
        "qrl" {
            .\Goshh-Client.exe -qrl
        }
        "qrc" {
            .\Goshh-Client.exe -qrc
        }
        "rune" {
            if ([string]::IsNullOrEmpty($in)) {
                Write-Error "User input is required for this command"
                return
            }
            .\Goshh-Client.exe -rune $in
        }
        "file" {
            if ([string]::IsNullOrEmpty($in)) {
                Write-Error "User input is required for this command"
                return
            }
            .\Goshh-Client.exe -file $in
        }
        "piped" {
            if ([string]::IsNullOrEmpty($in)) {
                Write-Error "User input is required for this command"
                return
            }
            "$in" | .\Goshh-Client.exe
        }
        "pqr" {
            if ([string]::IsNullOrEmpty($in)) {
                Write-Error "User input is required for this command"
                return
            }
            "$in" | .\Goshh-Client.exe -qr
        }
        "pqrl" {
            if ([string]::IsNullOrEmpty($in)) {
                Write-Error "User input is required for this command"
                return
            }
            "$in" | .\Goshh-Client.exe -qrl
        }
        "pqrc" {
            if ([string]::IsNullOrEmpty($in)) {
                Write-Error "User input is required for this command"
                return
            }
            "$in" | .\Goshh-Client.exe -qrc
        }
        Default {
            Write-Error "Invalid short call"
        }
    }
}
