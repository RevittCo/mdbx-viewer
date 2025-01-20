# MDBX Viewer

Welcome to **MDBX Viewer**! MDBX Viewer is a tool to help you visualise data in your MDBX database.

## Quick Start

1. **Pull the repo**

    ```bash
    git clone https://github.com/RevittConsulting/mdbx-viewer.git
    cd mdbx-viewer
    ```

2. **Add .env with your DATA_DIR**

    ```bash
    echo 'DATA_DIR="my/path"' > .env
    ```

3. **Run `make`**

    ```bash
    make
    ```

4. **Open `localhost:3000` in your web browser**

[http://localhost:3000](http://localhost:3000)

---

## Run on local dev mode

1. **Create .env in web**

You must create a .env and specify the API_URL. You can just copy the `.env.example` and rename to `.env`

2. **Start API server**

    ```bash
    cd api
    export DATA_DIR="your/path"
    go run cmd/api/main.go
    ```

2. **Start Frontend**

    ```bash
    cd web
    npm install
    npm run dev
    ```

2. **Open `localhost:5173` in your web browser**

[http://localhost:5173](http://localhost:5173)

---

Thank you for using MDBX Viewer!