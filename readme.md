# 🧱 Dockerize Monorepo Structure with Golang and TypeScript

Welcome to the **ultimate monorepo template** — designed for **high-performance services**, clean architecture, and **maximum developer experience**.

> ✨ Perfect for microservices, modular APIs, or scaling teams with shared utilities.

---

## ⚡️ Key Features

---

1. Dockerize Monorepo Structure could be implemented by containers, keep the local environment clean and consistent with the production environment.(to be implemented).

2. Turbo repo could be used to manage the monorepo structure, it could keep builded files as needed and hot reload.(to be implemented).

3. Dependency Injection and Factory Pattern could be used to manage the dependencies of the services.

4. Apps CLI tools shows each app's logs in terminal by selecting app.(to be implemented).

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

### GRPC generate

```bash
brew install bufbuild/buf/buf
pnpm setup
pnpm run buf:gen
```

2️⃣ Build All Packages

```
pnpm run build
```

3️⃣ Start API in Development

```bash
pnpm --filter @monorepo-services/api run start:dev
```

4️⃣ Start API in Production

```
pnpm --filter @monorepo-services/api run start:prod
```

## 💻 Contribution

Feel free to fork, improve, and submit PRs. Let’s make scalable backend monorepos easy for everyone 💪.
