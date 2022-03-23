import os

def main():
    os.system("go build gomaker.go")
    os.system("sudo cp gomaker /usr/local/bin/gomaker")
    print("Installation complete")

if __name__ == "__main__":
    main()
