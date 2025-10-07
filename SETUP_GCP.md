# 🔐 Cara Mendapatkan Google Cloud Credentials

Panduan lengkap untuk mendapatkan `GOOGLE_CREDENTIALS` (Service Account JSON) untuk deployment Discord bot ke Google Cloud Run.

---

## 📋 Prerequisites

- Google Account
- Google Cloud Platform access
- Credit card untuk enable billing (akan ada free tier $300)

---

## 🚀 Step-by-Step Guide

### Step 1: Buat/Pilih Google Cloud Project

1. Buka [Google Cloud Console](https://console.cloud.google.com/)
2. Login dengan Google Account Anda

3. **Buat Project Baru** atau pilih existing project:
   - Klik dropdown project di atas (next to "Google Cloud")
   - Klik "NEW PROJECT"
   - Isi:
     - **Project name**: `thomas-discord-bot` (atau nama lain)
     - **Project ID**: otomatis generate atau custom
     - **Location**: pilih organization jika ada
   - Klik "CREATE"
   - Tunggu beberapa detik, lalu pilih project yang baru dibuat

4. **Enable Billing** (wajib untuk Cloud Run):
   - Menu ☰ → "Billing"
   - Link billing account atau buat baru
   - Isi data credit card (akan ada free trial $300)

---

### Step 2: Enable Required APIs

Di Cloud Console, enable API berikut:

1. Menu ☰ → "APIs & Services" → "Library"
2. Search dan enable satu per satu:
   - ✅ **Cloud Run API**
   - ✅ **Cloud Build API**
   - ✅ **Artifact Registry API**
   - ✅ **Container Registry API**

Atau via Cloud Shell (klik icon >_ di atas):

```bash
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable artifactregistry.googleapis.com
gcloud services enable containerregistry.googleapis.com
```

---

### Step 3: Buat Service Account

1. Menu ☰ → "IAM & Admin" → "Service Accounts"
2. Klik "**+ CREATE SERVICE ACCOUNT**" di atas

#### 3.1 Service Account Details
- **Service account name**: `github-actions-bot`
- **Service account ID**: `github-actions-bot` (auto-fill)
- **Description**: `Service account for GitHub Actions deployment`
- Klik "**CREATE AND CONTINUE**"

#### 3.2 Grant Permissions (Roles)
Tambahkan roles berikut dengan klik "**+ ADD ANOTHER ROLE**":

| Role Name | Purpose |
|-----------|---------|
| **Cloud Run Admin** | Deploy & manage Cloud Run services |
| **Service Account User** | Act as service account |
| **Cloud Build Editor** | Build Docker images |
| **Artifact Registry Administrator** | Push/pull container images |
| **Storage Admin** | Access Cloud Storage |

Cara menambahkan:
- Klik dropdown "Select a role"
- Search role name (contoh: "Cloud Run Admin")
- Pilih role
- Klik "+ ADD ANOTHER ROLE" untuk tambah role lain
- Setelah semua role ditambahkan, klik "**CONTINUE**"

#### 3.3 Grant Users Access
- Skip bagian ini (opsional)
- Klik "**DONE**"

---

### Step 4: Generate JSON Key

1. Di halaman "Service Accounts", cari service account yang baru dibuat (`github-actions-bot`)
2. Klik **3 dots (⋮)** di kolom "Actions"
3. Pilih "**Manage keys**"
4. Klik "**ADD KEY**" → "**Create new key**"
5. Pilih format: **JSON**
6. Klik "**CREATE**"
7. File JSON akan otomatis ter-download ke komputer Anda

**⚠️ PENTING:**
- File ini berisi credentials yang sangat sensitif
- Jangan pernah commit ke Git/GitHub
- Simpan di tempat aman
- Jika hilang/ter-expose, delete key dan buat baru

---

### Step 5: Copy JSON Content

1. Buka file JSON yang ter-download dengan text editor
2. Isinya akan seperti ini:

```json
{
  "type": "service_account",
  "project_id": "your-project-id",
  "private_key_id": "abc123...",
  "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
  "client_email": "github-actions-bot@your-project-id.iam.gserviceaccount.com",
  "client_id": "123456789",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/..."
}
```

3. **Copy seluruh isi file** (termasuk { } dan semua karakter)

---

### Step 6: Add to GitHub Secrets

1. Buka GitHub repository: `thomas-bot-dc`
2. Pergi ke: **Settings** → **Secrets and variables** → **Actions**
3. Klik "**New repository secret**"
4. Isi:
   - **Name**: `GOOGLE_CREDENTIALS`
   - **Secret**: Paste seluruh isi JSON yang di-copy tadi
5. Klik "**Add secret**"

---

## ✅ Verification

Setelah semua setup, pastikan GitHub Secrets sudah lengkap:

| Secret Name | Status | Description |
|-------------|--------|-------------|
| `GOOGLE_CREDENTIALS` | ✅ | JSON Service Account |
| `DISCORD_TOKEN` | ⏳ | Token Discord (regenerate dulu!) |
| `CLIENT_ID` | ⏳ | Application ID Discord |
| `PUBLIC_KEY` | ⏳ | Public Key Discord |

---

## 🐛 Troubleshooting

### Error: "Permission denied" saat deploy

**Solusi**: Pastikan Service Account punya roles yang lengkap (lihat Step 3.2)

### Error: "API not enabled"

**Solusi**: Enable API yang diminta via Cloud Console atau command di Step 2

### JSON format error

**Solusi**:
- Pastikan copy JSON lengkap termasuk { } pembuka/penutup
- Jangan ada karakter tambahan
- Jangan edit isi JSON

### Service Account tidak muncul

**Solusi**: Pastikan sudah pilih project yang benar di Cloud Console

---

## 💰 Estimasi Cost

Dengan konfigurasi bot ini:
- **Memory**: 512Mi
- **CPU**: 1 core
- **Min instances**: 1 (always on)
- **Region**: asia-northeast1

Estimasi: **~$5-15/bulan**

**Tips hemat:**
- Ubah `min-instances: 0` di workflow (bot akan sleep saat tidak dipakai)
- Gunakan region yang lebih murah: `us-central1`
- Kurangi memory jadi `256Mi`

---

## 📚 Resources

- [GCP Service Accounts](https://cloud.google.com/iam/docs/service-accounts)
- [Cloud Run Pricing](https://cloud.google.com/run/pricing)
- [GitHub Actions Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)

---

🎉 **Setelah semua setup, Anda siap deploy bot ke GCP!**
