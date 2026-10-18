# zb-patch

ติดตั้ง **ZombieBuddy** (Java agent) ให้ Project Zomboid บน Windows อัตโนมัติ — รัน `zb-patch.exe` ครั้งเดียว (ขึ้น UAC ขอสิทธิแอดมินอัตโนมัติ) ไม่มี flag ไม่ถาม y/N

> ⚠ ต้อง [subscribe ZombieBuddy บน Steam](https://steamcommunity.com/sharedfiles/filedetails/?id=3619862853) และให้ Workshop โหลด mod ให้ครบก่อน

---

## ทำไมต้องใช้

Mod ZombieBuddy ต้องมี `ZombieBuddy.jar` + `zbNative.dll` ในโฟลเดอร์เกม และ `-agentlib:zbNative` ใน `ProjectZomboid64.json` ไม่งั้น Lua จะไม่มีตัวแปร `ZombieBuddy` (error แบบ `global is nil` / `setAutoFixModOrder of non-table: null` ตอน Mod Options)

---

## การทำงาน

`main.go` → **discover** (หา Steam / เกม / Workshop จาก Registry + `libraryfolders.vdf`) → **plan** (สรุปงาน) → **apply** (copy jar+dll · แก้ JSON · แก้ `ZombieBuddy_Options.lua` build 41/42 · backup `*.zb-patch.bak`)

```
รัน exe → แสดง path + รายการ → แก้ไฟล์ → Done. Restart Project Zomboid from Steam.
```

| โมดูล | ทำอะไร |
|--------|---------|
| `internal/discover` | หา path |
| `internal/plan` | เช็กว่าต้อง copy/แก้อะไร (ข้ามของที่ทำแล้ว) |
| `internal/launcher` | เพิ่ม `-agentlib:zbNative` ใน `vmArgs` |
| `internal/lua` | ใส่ `if not ZombieBuddy then return end` ใน `applySettings()` |

---

## ฟีเจอร์

- ✓ หา path เอง · ติดตั้ง agent + JSON + Lua guard ครบชุด · รันซ้ำได้
- ✗ Windows เท่านั้น · ไม่แก้ `ProjectZomboid64.bat` / Steam Launch Options ([installer ทางการ](https://github.com/zed-0xff/ZombieBuddy/releases/tag/windows_installer) ถ้าเปิดเกมแบบ BAT)

---

## ก่อนใช้

- Windows + Steam + Project Zomboid
- Subscribe mod ZombieBuddy แล้ว
- **Build เอง:** Go 1.22+ (`go version`)

---

## Build (ถ้ามีซอร์ส)

```powershell
cd patch
go mod tidy
go build -o zb-patch.exe .
```

---

## วิธีใช้

1. ปิดเกม (แนะนำ) → ดับเบิลคลิก **`zb-patch.exe`** (หรือ `.\zb-patch.exe`)
2. กด **Yes** ที่หน้าต่าง UAC (Administrator)
3. รอจน **`Done. Restart Project Zomboid from Steam.`** แล้วเปิดเกมจาก Steam ใหม่

---

## โครงโฟลเดอร์ (ย่อ)

```
patch/main.go          CLI
internal/discover      หา path
internal/plan          วางแผน
internal/apply         ลงมือแก้ + backup
internal/launcher      ProjectZomboid64.json
internal/lua           ZombieBuddy_Options.lua
```

---

## ประวัติการแก้ไข

### 1.3.0 — 2026-10-18

- **เพิ่ม:** ขอสิทธิ Administrator อัตโนมัติ (UAC) ตอนเปิดโปรแกรม — `internal/elevation`

### 1.2.0 — 2026-10-18

- **เปลี่ยน:** ย่อ README ให้สั้น อ่านเร็ว ครบใจความ (ตามสไตล์ createreadme 2.0)

### 1.1.0 — 2026-10-18

- **เปลี่ยน:** รัน exe แล้วทำงานทันที — ไม่มี flag / ไม่ถาม y/N

### 1.0.0 — 2026-10-18

- **เพิ่ม:** zb-patch — discover, copy agent, patch JSON/Lua, backup
