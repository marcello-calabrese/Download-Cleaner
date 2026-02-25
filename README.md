
# Download Cleaner

A simple Go program to help you clean up your Downloads folder by moving files to their appropriate locations based on file type. It categorizes files into folders like "Documents", "Images", "Videos", etc., and moves them accordingly.

## Features

- Automatically categorizes files based on their extensions.
- Provides a preview of the files that will be moved before executing the action.
- Supports a wide range of file types and can be easily extended to include more.

## Usage

1. Clone the repository:
   ```bash
   git clone
    ```
2. Navigate to the project directory:
   ```bash
    cd download-cleaner
    ``` 
3. Build the application:
   1. On Windows:
      ```bash
      go build -o download-cleaner.exe main.go
      ```
    2. On macOS/Linux:
        ```bash
        go build -o download-cleaner main.go
        ```
4. Run the application:
   On Windows:
      1. Double-click the `download-cleaner.exe` file in the project directory, or run it from the command line:
         ```bash

         .\download-cleaner.exe

## Troubleshooting

If you encounter any issues while running the application, please check the following:
- Ensure you have Go installed and properly set up on your system.
- Make sure you have the necessary permissions to access and modify files in your Downloads folder. 
- Check the console output for any error messages that may indicate what went wrong.
- If the application is not categorizing files correctly, verify that the file extensions are included in the categorization logic in the code.
- If you are running the application on Windows, ensure that you have the correct path to your Downloads folder set in the code.
- If you are still having trouble, feel free to open an issue on the GitHub repository with details about the problem you are facing.
- If you are running the application on macOS/Linux, ensure that you have the correct path to your Downloads folder set in the code.
- If you are running the application on Windows, ensure that you have the correct path to your Downloads folder set in the code.

## Contributing

Contributions are welcome! If you have any ideas for improvements or new features, please feel free to submit a pull request or open an issue on the GitHub repository.

## Bug Reports

If you encounter any bugs while using the application, please report them by opening an issue on the GitHub repository. Include as much detail as possible about the bug, including steps to reproduce it and any error messages you received.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.


## Acknowledgements

- [Go Programming Language](https://golang.org/) - The language used to build this application.

