<div align="center">
  <img src="assets/logo.jpg" alt="Spring Logo" width="200"/>
  <h1>⚡ Spring-CLI</h1>
  <p><strong>Un générateur ultra-rapide pour Spring Boot, écrit en Go.</strong></p>
  
  [![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
  [![Spring Boot](https://img.shields.io/badge/Spring_Boot-3.X-6DB33F?style=flat-square&logo=spring)](https://spring.io/projects/spring-boot)
  [![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
</div>

---

## 🚀 Qu'est-ce que Spring-CLI ?

**Spring-CLI** est un outil en ligne de commande (CLI) développé en Go, conçu pour les développeurs Java/Spring Boot. Fini la copie de code répétitif et les tâches fastidieuses ! 

Contrairement aux générateurs basés sur la JVM, `spring-cli` compile en un **binaire natif**, ce qui signifie :
- **Zéro temps de démarrage** (pas de JVM à chauffer)
- **Aucune dépendance lourde** (pas besoin de Node.js, Ruby, ou Python)
- **Génération instantanée** de votre code boilerplate (Entités, Repositories, Services, Controllers, DTOs).

---

## 📦 Installation

### 1️⃣ Depuis les sources (Développement)
Assurez-vous d'avoir Go installé sur votre machine.

```bash
git clone https://github.com/TALLHAMADOU/spring-boot-cli.git
cd spring-boot-cli
go build -o spring-cli .
sudo mv spring-cli /usr/local/bin/
```

### 2️⃣ Depuis les Releases GitHub
Téléchargez simplement la dernière archive depuis l'onglet [Releases](../../releases) et placez le binaire dans votre `PATH`.

---

## 💡 Utilisation Rapide

Placez-vous dans votre dossier de travail et commencez la magie :

### 🏗️ 1. Initialiser un projet Spring Boot
Créez instantanément un squelette Maven ou Gradle prêt à l'emploi.

```bash
spring-cli install:project maven --name demo --package com.example.demo
```

### 🧬 2. Générer une Entité JPA
Générez une entité avec ses attributs en une seule ligne :

```bash
spring-cli make entity User --fields "firstName:String, lastName:String, age:int"
```
> 💡 *Astuce : L'outil injectera automatiquement Lombok et Spring Data JPA dans votre `pom.xml` ou `build.gradle` si nécessaire.*

### 🛠️ 3. Générer le Boilerplate (Repository, Service, Controller)
Générez toute la couche logique autour de votre entité en un clin d'œil :

```bash
# Génère l'interface JpaRepository
spring-cli make repository User

# Génère l'interface Service et son implémentation
spring-cli make service User --entity User

# Génère un Controller REST avec toutes les routes CRUD (GET, POST, PUT, DELETE)
spring-cli make controller User --entity User --crud
```

### 🧪 4. Générer des DTOs & Tests
```bash
# Génère un DTO avec son Mapper MapStruct
spring-cli make dto User --fields "firstName:String, lastName:String" --mapper

# Génère des tests basés sur les champs de l'entité
spring-cli make test service User
```

### 🐳 5. Générer une configuration Docker
Générez un `Dockerfile` optimisé et un `docker-compose.yml` incluant la base de données de votre choix :

```bash
spring-cli make docker --db postgres --port 8080 --jdk 17
```

### 🚨 6. Gestion Globale des Exceptions
```bash
# Génère un ControllerAdvice avec les gestionnaires d'erreurs standards (400, 404, etc.)
spring-cli make exception-handler
```

---

## 📋 Liste des Commandes

| Commande | Description | Flags Principaux |
|----------|-------------|------------------|
| `install:project` | Initialise un projet Maven ou Gradle. | `--name`, `--package` |
| `make entity` | Génère une Entité JPA. | `--fields`, `--lombok`, `--uuid`, `--validate`, `--has-many` |
| `make repository`| Génère une interface JpaRepository. | `--uuid` |
| `make service` | Génère un Service (Interface + Impl). | `--entity` |
| `make controller`| Génère un RestController. | `--crud`, `--entity`, `--validate` |
| `make dto` | Génère un DTO (Data Transfer Object). | `--fields`, `--mapper`, `--validate` |
| `make test` | Génère des squelettes de tests JUnit. | |
| `make docker` | Génère un Dockerfile et docker-compose.yml. | `--db`, `--port`, `--jdk` |
| `make exception-handler`| Génère un `@ControllerAdvice`. | |
| `version` | Gère le versioning sémantique. | `--bump`, `--set`, `--tag` |

---

## 🔄 Gestion des Versions (Versioning)

`spring-cli` inclut un gestionnaire de version intégré qui détecte et met à jour automatiquement la version dans votre fichier `pom.xml` ou `build.gradle`.

```bash
# Voir la version actuelle
spring-cli version

# Incrémenter la version mineure (ex: 1.0.0 -> 1.1.0)
spring-cli version --bump minor

# Définir une version spécifique
spring-cli version --set 2.0.0
```

---

## 🤝 Contribuer

Les contributions sont grandement appréciées ! Si vous avez une idée d'amélioration ou trouvez un bug :
1. Ouvrez une **Issue**.
2. **Forkez** le projet.
3. Créez une branche pour votre fonctionnalité (`git checkout -b feature/ma-super-feature`).
4. Commitez vos changements (`git commit -m 'feat: ajout d'une super feature'`).
5. Poussez la branche (`git push origin feature/ma-super-feature`).
6. Ouvrez une **Pull Request**.

---

## 📜 Licence

Ce projet est sous licence **MIT**. Voir le fichier [LICENSE](LICENSE) pour plus de détails.

<div align="center">
  <i>Développé avec ❤️ pour rendre la vie des devs Spring plus facile.</i>
</div>
