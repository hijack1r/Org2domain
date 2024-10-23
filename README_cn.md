`Org2domain`是一款主域名发现工具，可以帮助你快速找到目标对应的主域名。 
该工具非常适用于渗透测试学习者，当你的领导甩给你一份组织清单，你还在浏览器中手动搜索对应主域名么?现在`Org2domain`可以帮你自动化且批量化地完成这件事，此时你可以悠闲地喝口咖啡。本工具对中英文组织名均可识别，输入的组织名不受语言限制。


# 用法
```shell
.\Org2domain64.exe -h
```
这将显示该工具的帮助，以下是它支持的所有参数。
```text
Usage:
  ./Org2domain64.exe [flags]

Flags:
  -h, --help  Help for Org2domain64
  -f string   File path of target file
  -p string   If need proxy to access Google, please set the proxy address
  -o string   Exported file name
  
  Such as:
  ./Org2domain64.exe -f OrgEn.txt -p http://127.0.0.1:7890 -o OrgEn.csv
```
![run](./run.jpg "run")
![verify](./verify.jpg "verify")

# Declaration
由于是短时间熬夜所写，头脑昏昏，料想会有不少错误，欢迎指出。本工具仅供学习使用，如在使用本工具的过程中存在任何非法行为，您将自行承担所有后果。