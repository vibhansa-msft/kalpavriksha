# kalpavriksha

## Azure Test Data Generator

Generate test data directly in your storage account.

## Configuration

- --dirs n : Number of directories to be generated
- --files n : Number of files to be generated in each directory
- --size n : Size of each file in MBs
- --concurrency n : Number of files being uploaded in parallel
- --type [ZERO/RANDOM/FILE] : Type of data to be written in each file. 
 
       -- ZERO : File will be filled with zeros
       -- RANDOM : File will be filled with random data
       -- FILE : Use source file data padded with zeros

- --src-file \<path\> : File path to be used as source data when --type=FILE is set.
- --dst-path \<path\> : Path in the container where test data needs to be generated. By default it will be generated on container root.
- --acct-type \<type\> : As of now only Blob type is supported
- --md5 true|false : Compute and set MD5 Sum for each file uploaded to container.
- --tier \<tier\> : Tier to be set for each file uploaded to container.
- --delete true|false : Instead of generating the data, delete previously generated data using this tool

## Environment Variables

- AZURE_STORAGE_ACCOUNT : Storage account name
- AZURE_STORAGE_ACCESS_KEY : Storage account key
- AZURE_STORAGE_ACCOUNT_CONTAINER : Container name where generated data will be stored

## Example

- Before running this command set the environment variables mentioned above and provide your account credentials
- Below command generated 100 directories each filled with 100 files each and each file of size 5MB. All this data will be stored in "dir1" in the container and each file will have its md5sum set with blob being in "cool" tier.
    
        -- .\kalpavriksha.exe --dirs 100 --files 100 --size 5 --tier "cool" --type "random" --dst-path "dir1" --concurrency 10 --md5 true
