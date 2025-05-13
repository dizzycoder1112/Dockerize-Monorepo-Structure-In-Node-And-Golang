# 🧱 Dockerize Monorepo Structure with Golang and TypeScript

Welcome to the **ultimate monorepo template** — designed for **high-performance services**, clean architecture, and **maximum developer experience**.

> ✨ Perfect for microservices, modular APIs, or scaling teams with shared utilities.

## ⚡️ Key Features

1. Dockerize Monorepo Structure could be implemented by containers, keep the local environment clean and consistent with the production environment.

2. Watching the packages changes, it could keep builded files as needed and hot reload.

3. Dependency Injection and Factory Pattern could be used to manage the dependencies of the services, let the program more flexible and easy to test.

4. log-tool CLI tools shows each app's logs in terminal by selecting app.

## 📂 Project Structure

```
root/
├── ts-packages/
│ └── shared/
│   └── src/
│     ├──constants/
│     └── utils/
│ └── logger/
│   └── src/
│ └── db/
│   └── src/
│ └── grpc/
│   └── src/
├── go-packages/
│ └── grpc/
├── apps/
│ └── ts-restful-api/
│   ├── src/
│   ├── tsconfig.json
│   └── tsconfig.prod.json
├── tsconfig.json
├── buf.gen.yaml
├── buf.yaml
└── pnpm-workspace.yaml
```

---

## 🛠 Usage

### 1️⃣ Install Dependencies

```bash
pnpm install
```

### 2️⃣ Start in Development

⚠️caution: please ensure docker is installed and running.

📝description: this dev mode was powered by docker continaer, and watch the packages changes to rebuild and restart the services by turbo.

```bash
pnpm run start:dev
```

watch the logs by using log-tool

```bash
pnpm run log-tool
```
![Screenshot 2025-05-13 at 20 39 44](https://github.com/user-attachments/assets/00c495aa-d560-43f5-bdad-be9148a0c7ed)

### 3️⃣ Build All Packages

```
pnpm run build
```



## Else

### GRPC generate

```bash
brew install bufbuild/buf/buf
pnpm setup
pnpm run buf:gen
```

## 💻 Contribution

Feel free to fork, improve, and submit PRs. Let’s make scalable backend monorepos easy for everyone 💪.
