# TF_API_Requester
Terraform Enterprise and Cloud API Requester

# How to use
Make sure set environment **TF_TOKEN** with your API key. <br>
Make sure you have archive Terraform manifest to .tar.gz archive <br>
```
    tar -zcvf "yourTFmanifest.tar.gz" -C "yourTFmanifestFolder" .
```

Run this command to trigger Terraform server
```
    go run run_task.go yourTFmanifest.tar.gz
```