# 🔐 Cara Mendapatkan Google Credentials via Terminal

## Berdasarkan Screenshot Anda

Saya lihat Anda sudah punya Service Account: **`go-gcp-deployer@gwa-project-472118.iam.gserviceaccount.com`**

Service Account ini sudah punya roles yang cukup lengkap! Anda bisa pakai Service Account ini.

---

## 🚀 Method 1: Download JSON Key via gcloud CLI (RECOMMENDED)

### Step 1: Install Google Cloud SDK

**Windows:**
Download dari: https://cloud.google.com/sdk/docs/install#windows

Atau via PowerShell:
```powershell
(New-Object Net.WebClient).DownloadFile("https://dl.google.com/dl/cloudsdk/channels/rapid/GoogleCloudSDKInstaller.exe", "$env:Temp\GoogleCloudSDKInstaller.exe")
& $env:Temp\GoogleCloudSDKInstaller.exe
```

**Setelah install, restart terminal**

---

### Step 2: Login ke Google Cloud

```bash
gcloud auth login
```

Browser akan terbuka, login dengan akun `gilarwahibul@gmail.com`

---

### Step 3: Set Project

```bash
gcloud config set project gwa-project-472118
```

---

### Step 4: Verify Service Account Exists

```bash
gcloud iam service-accounts list
```

Output harus menampilkan:
```
EMAIL: go-gcp-deployer@gwa-project-472118.iam.gserviceaccount.com
NAME: go-gcp-deployer
```

---

### Step 5: Create & Download JSON Key

```bash
# Buat key baru dan save ke file
gcloud iam service-accounts keys create gcp-key.json \
  --iam-account=go-gcp-deployer@gwa-project-472118.iam.gserviceaccount.com
```

**Output:**
```
created key [abc123...] of type [json] as [gcp-key.json]
```

File `gcp-key.json` akan dibuat di current directory.

---

### Step 6: Lihat Isi JSON

**Windows (PowerShell):**
```powershell
Get-Content gcp-key.json
```

**Git Bash / WSL:**
```bash
cat gcp-key.json
```

---

### Step 7: Copy JSON Content

**Windows (Copy to Clipboard):**
```powershell
Get-Content gcp-key.json | Set-Clipboard
```

**Git Bash:**
```bash
cat gcp-key.json | clip
```

Sekarang JSON sudah ter-copy ke clipboard! Tinggal paste ke GitHub Secrets.

---

## 🚀 Method 2: Via Cloud Console (Quick)

Kalau gcloud CLI ribet, bisa lewat browser:

### Step 1: Buka Cloud Console

```
https://console.cloud.google.com/iam-admin/serviceaccounts?project=gwa-project-472118
```

### Step 2: Click Service Account

Klik email: `go-gcp-deployer@gwa-project-472118.iam.gserviceaccount.com`

### Step 3: Manage Keys

- Tab **KEYS**
- Click **ADD KEY** → **Create new key**
- Select **JSON**
- Click **CREATE**

File JSON akan ter-download otomatis.

---

## 📋 Add to GitHub Secrets

### Via GitHub Website:

1. Buka repository: https://github.com/YOUR_USERNAME/thomas-bot-dc
2. **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Isi:
   - **Name**: `GOOGLE_CREDENTIALS`
   - **Secret**: Paste JSON content (sudah di-copy dari step sebelumnya)
5. Click **Add secret**

### Via GitHub CLI (gh):

```bash
# Install gh CLI: https://cli.github.com/

# Login
gh auth login

# Add secret
gh secret set GOOGLE_CREDENTIALS < gcp-key.json
```

---

## ⚠️ Keamanan

Setelah upload ke GitHub Secrets:

```bash
# DELETE file JSON dari komputer (PENTING!)
rm gcp-key.json

# Atau di Windows PowerShell:
Remove-Item gcp-key.json
```

Jangan pernah commit file JSON ke Git!

---

## ✅ Verify Service Account Roles

Cek apakah Service Account punya roles yang cukup:

```bash
gcloud projects get-iam-policy gwa-project-472118 \
  --flatten="bindings[].members" \
  --format='table(bindings.role)' \
  --filter="bindings.members:go-gcp-deployer@gwa-project-472118.iam.gserviceaccount.com"
```

**Yang dibutuhkan:**
- ✅ Artifact Registry Administrator
- ✅ Cloud Functions Developer
- ✅ Cloud Run Admin
- ✅ Cloud Run Service Agent
- ✅ Logs Viewer
- ✅ Project IAM Admin
- ✅ Service Account Token Creator
- ✅ Service Account User
- ✅ Service Usage Admin
- ✅ Storage Admin

Dari screenshot, Service Account Anda sudah punya semua roles ini! 👍

---

## 🎯 Quick Commands Summary

```bash
# 1. Login
gcloud auth login

# 2. Set project
gcloud config set project gwa-project-472118

# 3. Create key
gcloud iam service-accounts keys create gcp-key.json \
  --iam-account=go-gcp-deployer@gwa-project-472118.iam.gserviceaccount.com

# 4. Copy to clipboard (Windows)
Get-Content gcp-key.json | Set-Clipboard

# 5. Or just display
cat gcp-key.json

# 6. Delete after upload to GitHub
rm gcp-key.json
```

---

## 🐛 Troubleshooting

### Error: "gcloud: command not found"

**Solusi**: Install Google Cloud SDK terlebih dahulu (lihat Step 1)

### Error: "Permission denied"

**Solusi**: Pastikan Anda login dengan akun `gilarwahibul@gmail.com` yang punya akses ke project

### Error: "Service account does not exist"

**Solusi**: Cek nama service account dengan:
```bash
gcloud iam service-accounts list
```

---

🎉 **Setelah dapat JSON, langsung add ke GitHub Secrets dan siap deploy!**
