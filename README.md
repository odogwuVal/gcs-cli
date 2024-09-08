# gcscli 
gcscli is a command-line interface (CLI) tool for interacting with Google Cloud Storage (GCS). With gcscli, you can upload files to GCS buckets, get object attributes, and manage files directly from your terminal. 

# Features
- Upload Files: Easily upload files to a GCS bucket.
- Progress Indicator: Track your upload progress in real-time.
- Bucket Management: List objects in your GCS buckets.
- Secure Authentication: Authenticate using service account keys.
- User Management: Restrict access to authorized users only.

# Installation

#### Download the prebuilt binary

- Macos
```
    curl -O https://storage.googleapis.com/zimvest_mobile_application/gcscli/gcscli
    chmod +x gcscli
```

- Windows
```
    curl -O https://storage.googleapis.com/zimvest_mobile_application/gcscli/gcscli.exe
```


- Move the binary to your executables directory:
```
    sudo mv gcscli /usr/local/bin/
    chmod +x /usr/local/bin/gcscli
```
NOTE: use your operating system's guideline to move binary into executable path.

### Usage
Uploading Files
To upload a file to a GCS bucket, run:
```
    gcscli upload -p <object path in bucket> -o <object name in bucket> -f <path to local file> <bucket name>
```
For example:
```
gcscli upload -p folder/in/bucket -o myfile.txt -f ./myfile.txt my-bucket-name
```

Get Object
To get objects in a GCS bucket, use:
```
    gcscli get -o folder/in/bucket/file.txt -b <bucket name>(defaults to zimvest prod bucket)
```
## Authentication
The CLI uses environment variables for managing authorized users. Ensure you set the following:
GCSCLI_USER_TOKEN=<user-token>

Example(Macos)
```
export GCSCLI_USER_TOKEN=yourtoken
```


# Contribute
TODO: Explain how other users and developers can contribute to make your code better. 

If you want to learn more about creating good readme files then refer the following [guidelines](https://docs.microsoft.com/en-us/azure/devops/repos/git/create-a-readme?view=azure-devops). You can also seek inspiration from the below readme files:
- [ASP.NET Core](https://github.com/aspnet/Home)
- [Visual Studio Code](https://github.com/Microsoft/vscode)
- [Chakra Core](https://github.com/Microsoft/ChakraCore)