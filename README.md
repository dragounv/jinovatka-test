# Jinovatka

## Directory structure

### main.go

The entrypoint

### paths.go

Router setup and handler initialization

### Handlers

HTTP handlers

### Components

Templ components for rendering HTML

### Services

Application logic

### DB

Persistence

## Endpoints

### GET /

Index / Landing page

- basic info
- form to submit URLs

### GET /admin/

Main admin page.

- search the entire database

TODO:

- [ ] Redirect to page showning progress of the URLs
    - [x] Create new model for Seeds Group
    - [x] Save the Group to storage
    - [x] Create new group view
    - [x] Query the storage for list of seeds
    - [ ] Create seed detail view
    - [ ] Query the storage for seed details
- [ ] Refactor services to properly implement the layer model (add service interfaces)
- [ ] Think of a way to do proper configuration

- Možnost nechat uživatele pojmenovat skupinu?