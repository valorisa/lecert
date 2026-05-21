# Lecert

[![CI](https://github.com/valorisa/lecert/actions/workflows/ci.yml/badge.svg)](https://github.com/valorisa/lecert/actions/workflows/ci.yml)
[![Lint](https://github.com/valorisa/lecert/actions/workflows/lint.yml/badge.svg)](https://github.com/valorisa/lecert/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/valorisa/lecert)](https://goreportcard.com/report/github.com/valorisa/lecert)
[![Go Version](https://img.shields.io/github/go-mod/go-version/valorisa/lecert)](https://go.dev/)
[![License](https://img.shields.io/github/license/valorisa/lecert)](https://github.com/valorisa/lecert/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/valorisa/lecert?include_prereleases)](https://github.com/valorisa/lecert/releases)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-brightgreen)](https://github.com/valorisa/lecert#installation)

[🇬🇧 English](README.md)

Un outil CLI multiplateforme qui simplifie la gestion des certificats Let's Encrypt pour les utilisateurs de tous niveaux.
Lecert adapte son interface à trois modes d'interaction (débutant, standard, expert) afin que les nouveaux utilisateurs comme les administrateurs systèmes expérimentés puissent obtenir, renouveler et révoquer des certificats TLS sans friction.

## C'est Quoi et Pourquoi J'en Ai Besoin ?

### Le Problème en Langage Clair

Lorsque vous visitez un site web et voyez l'icône de cadenas dans la barre d'adresse de votre navigateur, cela signifie que la connexion entre votre ordinateur et le site web est chiffrée. Personne ne peut espionner ce que vous tapez (mots de passe, numéros de carte bancaire) ni altérer le contenu de la page. Cette protection est fournie par un **certificat TLS** — un petit fichier numérique qui prouve que le site web est bien celui qu'il prétend être.

Si vous hébergez votre propre site web, blog, API ou tout service accessible sur Internet, vous avez besoin d'un certificat TLS. Sans lui, les navigateurs afficheront un avertissement effrayant "Non sécurisé" à vos visiteurs, les moteurs de recherche vous classeront plus bas, et les données sensibles circuleront en clair.

### Qu'est-ce que Let's Encrypt ?

[Let's Encrypt](https://letsencrypt.org/) est une Autorité de Certification gratuite et à but non lucratif qui délivre des certificats TLS sans frais. Avant l'existence de Let's Encrypt (2015), les certificats coûtaient de l'argent et nécessitaient des démarches administratives manuelles. Désormais, tout le monde peut en obtenir gratuitement — mais le processus implique toujours des étapes techniques qui peuvent être intimidantes.

### Que Fait Lecert ?

Lecert est un outil en ligne de commande qui gère pour vous l'ensemble du cycle de vie des certificats :

1. **Obtenir** — Prouve que vous possédez un domaine et obtient un certificat de Let's Encrypt
2. **Renouveler** — Rafraîchit automatiquement votre certificat avant son expiration (ils durent 90 jours)
3. **Révoquer** — Invalide un certificat si votre serveur est compromis ou si vous perdez le contrôle du domaine

### Ai-je Besoin de Ceci ?

Vous avez besoin de Lecert (ou d'un outil similaire) si l'un de ces cas s'applique :

- Vous hébergez un site web, une application ou une API sur votre propre serveur
- Vous voyez des avertissements "Non sécurisé" en visitant votre site
- Votre certificat actuel a expiré et vous en avez besoin d'un nouveau
- Vous voulez HTTPS mais ne voulez pas payer pour des certificats
- Vous gérez plusieurs domaines et voulez un seul outil pour tous les gérer

Vous n'en avez PAS besoin si :

- Vous utilisez un hébergeur qui gère automatiquement les certificats (Vercel, Netlify, Heroku)
- Vous utilisez déjà Caddy (qui a HTTPS automatique intégré)
- Votre site est purement local/interne sans exposition sur Internet

### Comment Fonctionne la Vérification du Domaine ?

Let's Encrypt a besoin d'une preuve que vous possédez effectivement le domaine avant de délivrer un certificat. Cela empêche quelqu'un d'obtenir un certificat pour un domaine qu'il ne contrôle pas. Il existe deux méthodes courantes :

**Challenge HTTP (recommandé pour les débutants) :** Let's Encrypt demande à votre serveur de placer un fichier spécifique à une URL spécifique. Si votre serveur peut répondre correctement, cela prouve que vous contrôlez le domaine. Cela nécessite que le port 80 soit ouvert sur votre serveur.

**Challenge DNS (pour les configurations avancées) :** Au lieu de placer un fichier sur votre serveur, vous ajoutez un enregistrement DNS spécifique à la configuration de votre domaine. C'est utile lorsque le port 80 est bloqué, lorsque vous utilisez un CDN, ou lorsque vous avez besoin de certificats pour des serveurs internes qui ne sont pas directement accessibles depuis Internet.

### Choisir Votre Niveau de Confort

Lecert propose trois modes car les gens ont différents niveaux d'expérience :

| Si vous êtes... | Utilisez ce mode | À quoi ça ressemble |
|-----------------|------------------|---------------------|
| Nouveau dans les serveurs et certificats | `--mode novice` | Un assistant amical qui pose 3 questions simples |
| À l'aise avec les outils en ligne de commande | `--mode standard` (par défaut) | Interface familière basée sur des flags comme les autres outils CLI |
| Un administrateur système qui veut un contrôle total | Commande `obtain-expert` | Toutes les options du protocole ACME exposées, rien de caché |

Vous pouvez toujours commencer avec le mode débutant et passer au mode standard ou expert au fur et à mesure que vous gagnez en confiance. Les certificats produits sont identiques quel que soit le mode utilisé.

## Table des Matières

- [C'est Quoi et Pourquoi J'en Ai Besoin ?](#cest-quoi-et-pourquoi-jen-ai-besoin-)
- [Fonctionnalités](#fonctionnalités)
- [Installation](#installation)
- [Démarrage Rapide](#démarrage-rapide)
- [Modes d'Utilisation](#modes-dutilisation)
- [Commandes](#commandes)
- [Fournisseurs DNS](#fournisseurs-dns)
- [Renouvellement Automatique](#renouvellement-automatique)
- [Configuration](#configuration)
- [Architecture](#architecture)
- [Développement](#développement)
- [Tests](#tests)
- [Contribuer](#contribuer)
- [Questions Fréquemment Posées](#questions-fréquemment-posées)
- [Glossaire](#glossaire)
- [Signalement de Vulnérabilités de Sécurité](#signalement-de-vulnérabilités-de-sécurité)
- [Licence](#licence)
- [Journal des Modifications](#journal-des-modifications)

## Fonctionnalités

- **Trois modes d'interaction** qui s'adaptent à l'expertise de l'utilisateur (assistant débutant, CLI standard, mode expert)
- **Multiplateforme** : binaire unique pour Linux, macOS et Windows sans aucune dépendance d'exécution
- **Challenges HTTP-01 et DNS-01** avec support intégré pour Cloudflare, Route53 et DigitalOcean
- **Planification automatique du renouvellement** via cron (Linux/macOS) ou Planificateur de Tâches (Windows)
- **Stockage sécurisé des clés** avec permissions POSIX (0600) sur Unix et restrictions ACL sur Windows
- **Renouvellement par lot** de tous les certificats gérés approchant l'expiration (seuil de 30 jours)
- **Inventaire des certificats** avec affichage du statut (valide, expirant bientôt, expiré)
- **Support de l'environnement de staging Let's Encrypt** pour les tests sans toucher aux limites de taux de production
- **Intégration GoReleaser** pour des builds de release reproductibles et multi-compilés
- **Détection automatique du fournisseur DNS** basée sur l'environnement pour des challenges DNS sans configuration

## Installation

### Depuis les Binaires de Release

Téléchargez le dernier binaire pour votre plateforme depuis la page [Releases](https://github.com/valorisa/lecert/releases).

```bash
# Linux (amd64)
curl -LO https://github.com/valorisa/lecert/releases/latest/download/lecert_linux_amd64.tar.gz
tar xzf lecert_linux_amd64.tar.gz
sudo mv lecert /usr/local/bin/

# macOS (Apple Silicon)
curl -LO https://github.com/valorisa/lecert/releases/latest/download/lecert_darwin_arm64.tar.gz
tar xzf lecert_darwin_arm64.tar.gz
sudo mv lecert /usr/local/bin/

# Windows (amd64) — extrayez le zip et ajoutez au PATH
```

### Depuis les Sources

Nécessite Go 1.21 ou supérieur.

```bash
go install github.com/valorisa/lecert/cmd/lecert@latest
```

### Compilation depuis le Dépôt

```bash
git clone https://github.com/valorisa/lecert.git
cd lecert
make build
# Binaire disponible dans ./bin/lecert
```

## Démarrage Rapide

### Pour les Débutants (Mode Novice)

L'assistant pose exactement trois questions et gère tout le reste automatiquement.

```bash
lecert --mode novice cert obtain
```

```text
=== Assistant de Certificat Let's Encrypt ===

1/3 Nom de domaine (ex : exemple.com) : monsite.com
2/3 Votre email (pour les alertes d'expiration de certificat) : moi@monsite.com
3/3 Comment devons-nous vérifier la propriété du domaine ?
    [1] HTTP (nécessite le port 80 ouvert) — recommandé pour la plupart des configurations
    [2] DNS  (nécessite un accès API au fournisseur DNS)
    Choix [1] : 1

Compris ! Demande de certificat pour monsite.com via challenge http-01...

Certificat obtenu pour monsite.com
  Expire : 2026-08-17
  Stocké : /home/utilisateur/.lecert/certs/monsite.com
```

### Pour les Utilisateurs Standard

```bash
lecert cert obtain --domain monsite.com --email moi@monsite.com --staging
```

### Pour les Experts

```bash
lecert cert obtain-expert \
  --domain monsite.com \
  --domain www.monsite.com \
  --email moi@monsite.com \
  --challenge dns-01 \
  --dns-provider cloudflare \
  --key-type ec384 \
  --preferred-chain "ISRG Root X1" \
  --timeout 5m
```

## Modes d'Utilisation

| Mode | Flag | Public | Comportement |
|------|------|--------|--------------|
| Novice | `--mode novice` | Utilisateurs débutants | Assistant interactif, maximum 3 questions, valeurs par défaut sensées |
| Standard | `--mode standard` | Utilisateurs réguliers | Basé sur des flags, messages d'erreur clairs, flags requis appliqués |
| Expert | Commande `obtain-expert` | Administrateurs système | Toutes les options ACME exposées, SAN multi-domaines, sélection du type de clé |

Le mode est défini globalement via le flag `--mode` et s'applique à toutes les sous-commandes.

## Commandes

### Opérations sur les Certificats

| Commande | Description |
|----------|-------------|
| `lecert cert obtain` | Obtenir un nouveau certificat (respecte le mode actuel) |
| `lecert cert obtain-expert` | Obtenir avec contrôle ACME complet (expert uniquement) |
| `lecert cert renew --domain X` | Renouveler un certificat spécifique |
| `lecert cert renew-all` | Renouveler tous les certificats expirant dans les 30 jours |
| `lecert cert revoke --domain X` | Révoquer un certificat (avec confirmation interactive) |
| `lecert cert list` | Afficher tous les certificats gérés avec leur statut |

### Planification

| Commande | Description |
|----------|-------------|
| `lecert schedule install` | Installer le renouvellement automatique (cron ou Planificateur de Tâches) |
| `lecert schedule uninstall` | Supprimer la planification du renouvellement automatique |
| `lecert schedule status` | Vérifier si le renouvellement automatique est actif |

### Flags Globaux

| Flag | Défaut | Description |
|------|--------|-------------|
| `--mode` | `standard` | Mode d'interaction : novice, standard, expert |
| `--version` | | Afficher la version et quitter |
| `--help` | | Afficher l'aide pour n'importe quelle commande |

## Fournisseurs DNS

Lecert supporte les challenges DNS-01 via les fournisseurs suivants. L'authentification est configurée via des variables d'environnement.

| Fournisseur | Variables d'Environnement |
|-------------|---------------------------|
| Cloudflare | `CF_DNS_API_TOKEN` ou `CF_API_EMAIL` + `CF_API_KEY` |
| AWS Route53 | `AWS_ACCESS_KEY_ID` + `AWS_SECRET_ACCESS_KEY` + `AWS_REGION` |
| DigitalOcean | `DO_AUTH_TOKEN` |

Lecert détecte automatiquement votre fournisseur DNS depuis les variables d'environnement. Vous pouvez aussi le spécifier explicitement :

```bash
export CF_DNS_API_TOKEN="votre-token-ici"
lecert cert obtain --domain monsite.com --email moi@monsite.com --challenge dns-01 --dns-provider cloudflare
```

## Renouvellement Automatique

Lecert peut installer une tâche planifiée qui renouvelle automatiquement les certificats approchant l'expiration.

### Installation

```bash
# Vérification quotidienne à 02:30 (par défaut)
lecert schedule install

# Intervalle personnalisé
lecert schedule install --interval twice-daily
lecert schedule install --interval hourly
```

### Comment Ça Marche

Le planificateur exécute `lecert cert renew-all --quiet` à l'intervalle configuré. Cette commande itère sur tous les certificats gérés et renouvelle ceux qui expirent dans les 30 jours. Les échecs de renouvellement sont enregistrés dans stderr mais n'arrêtent pas le processus pour les certificats restants.

### Détails par Plateforme

| Plateforme | Mécanisme | Emplacement de la Planification |
|------------|-----------|--------------------------------|
| Linux | Crontab utilisateur | `crontab -l` |
| macOS | Crontab utilisateur | `crontab -l` |
| Windows | Planificateur de Tâches | `schtasks /Query /TN LecertAutoRenew` |

## Configuration

### Emplacement de Stockage

Les certificats et métadonnées sont stockés dans `~/.lecert/certs/` par défaut. Remplacer avec la variable d'environnement `LECERT_DIR`.

```bash
export LECERT_DIR=/etc/lecert/certs
lecert cert obtain --domain monsite.com --email moi@monsite.com
```

### Structure des Répertoires

```text
~/.lecert/certs/
└── monsite.com/
    ├── cert.pem      # Chaîne de certificats (0644)
    ├── key.pem       # Clé privée (0600)
    └── meta.json     # Métadonnées (domaine, email, challenge, expiration)
```

### Sécurité

Les clés privées sont stockées avec des permissions restrictives (0600 sur Unix). Sur Windows, des ACL standard réservées à l'utilisateur s'appliquent. Les clés privées ne quittent jamais le système de fichiers local et ne sont jamais enregistrées dans les logs ni transmises.

## Architecture

```text
lecert/
├── cmd/lecert/              # Point d'entrée CLI et définitions des commandes
│   ├── main.go              # Commande racine, flag de mode, version
│   ├── obtain.go            # Obtention standard + routage novice
│   ├── obtain_expert.go     # Mode expert avec tous les flags ACME
│   ├── renew.go             # Renouvellement d'un seul domaine
│   ├── renew_all.go         # Renouvellement par lot (seuil J-30)
│   ├── revoke.go            # Révocation avec confirmation
│   ├── list.go              # Affichage de l'inventaire des certificats
│   └── schedule.go          # Installation/désinstallation/statut du planificateur
├── internal/
│   ├── acme/                # Wrapper du client ACME autour de lego
│   │   ├── acme.go          # Opérations Obtain, Renew, Revoke
│   │   └── dns.go           # Registre et détection des fournisseurs DNS
│   ├── store/               # Stockage sécurisé des certificats
│   │   └── store.go         # Save, Load, List avec permissions
│   ├── wizard/              # Assistant interactif mode novice
│   │   └── wizard.go        # Flux guidé en 3 questions
│   └── scheduler/           # Planification du renouvellement automatique
│       ├── scheduler.go     # Dispatcher agnostique de l'OS
│       ├── cron.go          # Gestion de crontab Linux/macOS
│       └── windows.go       # Gestion du Planificateur de Tâches Windows
├── Makefile                 # Cibles build, test, release
├── .goreleaser.yaml         # Configuration de release
├── go.mod
└── go.sum
```

### Décisions de Conception

- **Bibliothèque lego** : Fournit un support mature du protocole ACME avec 100+ fournisseurs DNS, évitant la réinvention
- **Framework CLI cobra** : Standard de l'industrie pour les CLI Go, permet les sous-commandes et l'auto-complétion
- **Packages internal/** : Empêchent les imports externes, appliquent les limites d'API
- **Pas de CGO** : Permet une vraie compilation croisée sans dépendances de chaîne d'outils
- **Fichiers de métadonnées (meta.json)** : Permettent un renouvellement sans état en persistant l'email et le type de challenge par domaine

## Développement

### Prérequis

- Go 1.21 ou supérieur
- Make (optionnel, pour les cibles de commodité)
- GoReleaser (optionnel, pour les builds de release)

### Compilation

```bash
make build          # Compiler pour la plateforme actuelle
make release        # Compilation croisée toutes plateformes
make test           # Exécuter tous les tests
make clean          # Supprimer les artefacts de build
```

### Cibles de Compilation Croisée

| OS | Architecture | Nom du Binaire |
|----|--------------|----------------|
| Linux | amd64 | `lecert-linux-amd64` |
| Linux | arm64 | `lecert-linux-arm64` |
| macOS | amd64 | `lecert-darwin-amd64` |
| macOS | arm64 | `lecert-darwin-arm64` |
| Windows | amd64 | `lecert-windows-amd64.exe` |

## Tests

```bash
# Exécuter tous les tests
go test ./... -v

# Exécuter les tests d'un package spécifique
go test ./internal/store/ -v
go test ./internal/wizard/ -v
go test ./internal/acme/ -v

# Exécuter avec couverture
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Couverture de Test

| Package | Tests | Couverture |
|---------|-------|-----------|
| internal/store | 3 tests (Save/Load, List, Permissions) | Opérations CRUD de base |
| internal/wizard | 4 tests (entrée valide, choix DNS, domaine vide, EOF) | Tous les chemins de l'assistant |
| internal/acme | 5 tests (détection DNS, registre des fournisseurs) | Logique des fournisseurs |

### Environnement de Staging

Pour les tests d'intégration sans toucher aux limites de taux de production de Let's Encrypt, utilisez le flag `--staging` sur toutes les commandes. Les certificats de staging ne sont pas approuvés par les navigateurs mais exercent le flux ACME complet.

## Contribuer

Veuillez lire [CONTRIBUTING.md](CONTRIBUTING.md) pour les détails sur notre code de conduite, le flux de travail de développement et le processus de soumission des pull requests.

### Guide Rapide de Contribution

1. Forkez le dépôt
2. Créez votre branche de fonctionnalité (`git checkout -b feature/fonctionnalite-incroyable`)
3. Exécutez les tests (`make test`)
4. Commitez vos modifications en suivant [Conventional Commits](https://www.conventionalcommits.org/)
5. Poussez vers la branche (`git push origin feature/fonctionnalite-incroyable`)
6. Ouvrez une Pull Request

## Questions Fréquemment Posées

### Je suis un débutant complet. Est-ce que cela va casser mon serveur ?

Non. Lecert ne modifie jamais la configuration de votre serveur web (Nginx, Apache, etc.). Il obtient uniquement des fichiers de certificat et les stocke dans un dossier. Vous devez toujours configurer votre serveur web pour utiliser ces fichiers, mais Lecert ne peut pas accidentellement casser quoi que ce soit qui fonctionne déjà.

Si vous utilisez le mode `--staging`, les certificats produits ne sont pas approuvés par les navigateurs mais sont par ailleurs identiques — cela vous permet de pratiquer l'ensemble du flux sans aucun risque ni limitation de taux.

### Que se passe-t-il si mon certificat expire ?

Les certificats Let's Encrypt sont valides 90 jours. La commande `schedule install` de Lecert configure un renouvellement automatique qui vérifie quotidiennement et renouvelle tout certificat expirant dans les 30 jours. Si vous oubliez de configurer le renouvellement automatique, votre site affichera un avertissement de navigateur après 90 jours, mais rien n'est cassé de façon permanente — exécutez simplement `lecert cert renew --domain votredomaine.com` pour le corriger.

### Dois-je être root/administrateur ?

Généralement non. Lecert stocke les certificats dans votre répertoire personnel (`~/.lecert/certs/`) et utilise un serveur HTTP non privilégié sur le port 5002 pour les challenges. Le seul cas où root pourrait être nécessaire est si vous voulez utiliser directement le port 80 pour les challenges HTTP (les ports en dessous de 1024 nécessitent des privilèges élevés sur la plupart des systèmes). L'approche recommandée est d'utiliser un reverse proxy ou une redirection de port à la place.

### Puis-je utiliser ceci avec Nginx/Apache/Caddy ?

Oui. Lecert produit des fichiers PEM standard (`cert.pem` et `key.pem`). Pointez la configuration de votre serveur web vers ces fichiers :

```nginx
# Exemple Nginx
ssl_certificate     /home/vous/.lecert/certs/votredomaine.com/cert.pem;
ssl_certificate_key /home/vous/.lecert/certs/votredomaine.com/key.pem;
```

```apache
# Exemple Apache
SSLCertificateFile    /home/vous/.lecert/certs/votredomaine.com/cert.pem
SSLCertificateKeyFile /home/vous/.lecert/certs/votredomaine.com/key.pem
```

Après le renouvellement, rechargez votre serveur web pour utiliser le nouveau certificat (`systemctl reload nginx`).

### Quelle est la différence entre les challenges HTTP et DNS ?

| Aspect | Challenge HTTP | Challenge DNS |
|--------|---------------|---------------|
| Nécessite | Port 80 ouvert sur votre serveur | Accès à l'API de votre fournisseur DNS |
| Fonctionne pour | Serveurs directement accessibles depuis Internet | N'importe quel domaine, même derrière des pare-feux |
| Difficulté | Plus facile (juste ouvrir un port) | Légèrement plus complexe (nécessite la configuration d'un token API) |
| Meilleur pour | Configurations simples sur serveur unique | Utilisateurs de CDN, serveurs internes, architectures complexes |

Si vous n'êtes pas sûr, commencez avec HTTP. Vous pouvez toujours passer à DNS plus tard.

### J'ai eu une erreur de limite de taux. Que faire ?

Let's Encrypt limite la délivrance de certificats à 5 par domaine par semaine en production. Si vous testez ou apprenez, utilisez toujours `--staging` pour éviter d'atteindre ces limites. Le staging a des limites beaucoup plus élevées et est conçu pour l'expérimentation. Une fois que votre configuration fonctionne avec staging, retirez le flag `--staging` pour obtenir un vrai certificat.

### En quoi est-ce différent de Certbot ?

Certbot est le client officiel Let's Encrypt maintenu par l'EFF. Il est excellent mais orienté vers les administrateurs système. Lecert diffère de plusieurs façons :

- **Binaire unique** — pas de dépendances Python, pas de pip, pas de virtualenv
- **Mode novice** — un assistant guidé qui pose 3 questions au lieu de nécessiter la connaissance des flags
- **Multiplateforme** — le même binaire fonctionne sur Linux, macOS et Windows
- **Ne touche pas à votre serveur web** — Certbot peut auto-configurer Nginx/Apache, ce qui est pratique mais peut aussi casser les configurations. Lecert produit uniquement des fichiers.
- **Portée plus légère** — Lecert fait une chose (cycle de vie des certificats) et vous laisse la configuration du serveur

Les deux outils produisent les mêmes certificats depuis la même infrastructure Let's Encrypt.

## Glossaire

Termes que vous pourriez rencontrer en travaillant avec des certificats :

| Terme | Signification |
|-------|---------------|
| **TLS** | Transport Layer Security — le protocole qui chiffre le trafic web (successeur de SSL) |
| **Certificat** | Un fichier numérique qui prouve l'identité d'un serveur et permet des connexions chiffrées |
| **Clé privée** | Un fichier secret que seul votre serveur devrait avoir — ne le partagez jamais |
| **Autorité de Certification (CA)** | Une organisation de confiance pour délivrer des certificats (Let's Encrypt est une CA) |
| **ACME** | Automatic Certificate Management Environment — le protocole utilisé par Let's Encrypt |
| **Domaine** | L'adresse de votre site web (ex : `exemple.com`) |
| **Challenge** | Un test qui prouve que vous contrôlez un domaine avant qu'un certificat ne soit délivré |
| **PEM** | Un format de fichier pour les certificats et clés (les fichiers `.pem` produits par Lecert) |
| **Renouvellement** | Obtenir un certificat frais avant que le certificat actuel n'expire |
| **Révocation** | Invalider un certificat (ex : si votre serveur a été compromis) |
| **Staging** | Un environnement de test qui délivre de faux certificats pour la pratique |

## Signalement de Vulnérabilités de Sécurité

Si vous découvrez une vulnérabilité de sécurité, veuillez suivre le processus de divulgation responsable décrit dans [SECURITY.md](SECURITY.md). N'ouvrez pas d'issue publique pour les vulnérabilités de sécurité.

## Licence

Ce projet est sous licence MIT. Voir le fichier [LICENSE](LICENSE) pour les détails.

## Journal des Modifications

Voir [CHANGELOG.md](CHANGELOG.md) pour une liste des modifications dans chaque release.
