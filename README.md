# 🎮 Power 4

Une implémentation web moderne du **Puissance 4**, développée en **Go** avec une interface neon / glassmorphism et des animations fluides.

Trois niveaux de difficulté, une interface responsive et un serveur Go minimaliste sans dépendance externe.

---

## 📖 Fonctionnalités principales

### 🕹️ Gameplay

* **3 niveaux** :

  * 🟢 **Facile** (6×7)
  * 🟡 **Normal** (6×9)
  * 🔴 **Difficile** (7×8)
* Détection automatique de **victoire** et de **match nul**
* Validation des coups et compteur de tours
* Mise en surbrillance du **dernier coup**

### 💻 Interface

* Design **neon / glassmorphism**
* Animations et mises à jour asynchrones (Fetch API)
* **Responsive** : mobile, tablette et desktop
* **Accessibilité** : ARIA, `prefers-reduced-motion`

### ⚙️ Technique

* Serveur **Go pur** (`net/http`, `html/template`)
* Logique de jeu optimisée et gestion d’erreurs


---

## 🚀 Installation rapide

```bash
git clone https://github.com/votre-username/power4.git
cd power4
go run .
```

➡️ Puis ouvrez : `http://localhost:8080`

---

## 🧠 Structure du projet

```
power4/
├── main.go         # Serveur HTTP
├── game.go         # Logique du jeu
├── handlers.go     # Routes et gestion des requêtes
├── templates/      # HTML (index, start)
└── static/         # CSS & JS
```
---

## 👥 Auteurs

* **yan1s973** – Backend & logique de jeu
* **maxence-ynov** – Frontend & design

> Projet réalisé dans le cadre du cours **Info B1**.

---


