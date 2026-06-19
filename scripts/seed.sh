#!/bin/bash

BASE="http://localhost:3000"

echo "=== Creando usuarios ==="

ADMIN1=$(curl -s -X POST $BASE/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Verónica Zacarias", "email": "vero@test.com", "password": "secret123", "role": "admin"}')
echo "Admin 1: $ADMIN1" | jq .

ADMIN2=$(curl -s -X POST $BASE/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Carlos Ruiz", "email": "carlos@test.com", "password": "secret123", "role": "admin"}')
echo "Admin 2: $ADMIN2" | jq .

USER=$(curl -s -X POST $BASE/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Lector Perez", "email": "lector@test.com", "password": "secret123", "role": "user"}')
echo "Usuario normal: $USER" | jq .

echo ""
echo "=== Login como admin ==="

TOKEN=$(curl -s -X POST $BASE/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "vero@test.com", "password": "secret123"}' | jq -r '.access_token')
echo "Token obtenido: ${TOKEN:0:30}..."

echo ""
echo "=== Creando libros ==="

BOOK1=$(curl -s -X POST $BASE/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Clean Code", "author": "Robert Martin", "isbn": "9780132350884", "year": 2008, "genre": "tech", "available_copies": 3}')
echo "Libro 1: $BOOK1" | jq .

BOOK2=$(curl -s -X POST $BASE/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "The Pragmatic Programmer", "author": "David Thomas", "isbn": "9780135957059", "year": 2019, "genre": "tech", "available_copies": 2}')
echo "Libro 2: $BOOK2" | jq .

BOOK3=$(curl -s -X POST $BASE/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Sapiens", "author": "Yuval Noah Harari", "isbn": "9780062316097", "year": 2011, "genre": "history", "available_copies": 5}')
echo "Libro 3: $BOOK3" | jq .

BOOK4=$(curl -s -X POST $BASE/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "El nombre de la rosa", "author": "Umberto Eco", "isbn": "9788423309702", "year": 1980, "genre": "fiction", "available_copies": 4}')
echo "Libro 4: $BOOK4" | jq .

echo ""
echo "=== Creando préstamos ==="

LOAN1=$(curl -s -X POST $BASE/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"user_id": 1, "book_id": 1}')
echo "Préstamo 1 (user 1 - libro 1): $LOAN1" | jq .

LOAN2=$(curl -s -X POST $BASE/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"user_id": 2, "book_id": 2}')
echo "Préstamo 2 (user 2 - libro 2): $LOAN2" | jq .

LOAN3=$(curl -s -X POST $BASE/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"user_id": 3, "book_id": 3}')
echo "Préstamo 3 (user 3 - libro 3): $LOAN3" | jq .

echo ""
echo "=== Estado final ==="
echo "Libros:"
curl -s $BASE/books | jq '.data[] | {id, title, available_copies}'
echo "Préstamos activos user 1:"
curl -s $BASE/loans/users/1 -H "Authorization: Bearer $TOKEN" | jq .