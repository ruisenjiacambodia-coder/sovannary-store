# 🛍️ Sovannary Store — Digital Menu & Inventory

កម្មវិធីគ្រប់គ្រងស្តុក និងបង្ហាញមុខទំនិញបែបឌីជីថល សម្រាប់ហាង **Sovannary Store**។

## ✨ មុខងារសំខាន់ៗ

### 🛒 Customer View
- បង្ហាញផលិតផលជា Card Grid ស្អាតៗ
- ស្វែងរក និងតម្រៀបតាមប្រភេទ/តម្លៃ
- ប្រព័ន្ធកន្ត្រកទំនិញ (Shopping Cart)
- Checkout ផ្ញើទៅ **WhatsApp** ឬ **Telegram**

### ⚙️ Admin Dashboard
- បន្ថែម/កែសម្រួល/លុប ផលិតផល
- គ្រប់គ្រងស្តុក (+1 / -1 / +10 / -10)
- មើលការបញ្ជាទិញ
- កំណត់ Telegram Bot / WhatsApp
- របាយការណ៍ស្តុក

### 📱 PWA Features
- ដំណើរការ **Offline** បាន
- **Background Sync** ពេលមានអ៊ីនធឺណិតវិញ
- **IndexedDB** សម្រាប់រក្សាទុកទិន្នន័យ
- Install ជា App លើទូរស័ព្ទ

## 🚀 ការដំឡើង

### តម្រូវការ
- Go 1.22+ (https://go.dev/dl/)

### ដំណើរការ
```bash
# 1. Clone ឬ copy ឯកសារទាំង 4 ចូល folder
mkdir sovannary-store && cd sovannary-store
# (copy main.go, dashboard.html, sw.js, manifest.json)

# 2. Run server
go run main.go

# 3. បើក browser
# http://localhost:8080
```

### Build Binary (Production)
```bash
go build -o sovannary main.go
./sovannary
```

## 🔧 ការកំណត់ Telegram Bot

1. និយាយជាមួយ `@BotFather` លើ Telegram
2. បង្កើត bot ថ្មី → យក **Bot Token**
3. បន្ថែម bot ចូល group/channel របស់អ្នក
4. រក **Chat ID** (ប្រើ `@userinfobot` ឬ `https://api.telegram.org/bot<TOKEN>/getUpdates`)
5. បំពេញក្នុង **Admin → Settings**

## 📂 រចនាសម្ព័ន្ធ API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/products` | បញ្ជីផលិតផលទាំងអស់ |
| POST | `/api/products` | បន្ថែមផលិតផលថ្មី |
| PUT | `/api/products/:id` | កែសម្រួលផលិតផល |
| DELETE | `/api/products/:id` | លុបផលិតផល |
| GET | `/api/orders` | បញ្ជីការបញ្ជាទិញ |
| POST | `/api/orders` | បង្កើតការបញ្ជាទិញថ្មី |
| GET | `/api/settings` | មើលការកំណត់ |
| PUT | `/api/settings` | កែការកំណត់ |
| POST | `/api/sync` | Sync ទិន្នន័យពី offline |

## 🎨 Design System

- **Primary:** Zinc-950 (background), Gold-500 (accent)
- **Fonts:** Battambang (Khmer), Playfair Display (Display), Inter (Body)
- **Style:** Luxury, Minimal, Modern

## 📝 License

MIT © 2026 Sovannary Store