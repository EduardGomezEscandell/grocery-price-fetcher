# Grocery price fetcher

Aquest programa permet descarregar els preus de diversos productes de les webs de varis supermercats.
De moment permet Bonpreu i Mercandona, però és fàcilment extensible.

## Instal·lació
El primer cop que l'utilitzes l'hauràs de compilar

1. Instala les dependències:
   - Ubuntu
     ```
     sudo apt install -y make golang-go
     ```
   - Fedora
     ```
     sudo dnf install -y make golang-go
     ```
   - Windows i altres: Descarrega Go de la [web](https://go.dev/dl/).
2. Compila
   - Linux: `make build-go`
   - Windows: `New-Item -Path bin -Type Directory ; go build -o ./bin/compra.exe build/cmd/compra`

## Entrada de dades
Per poder buscar els productes, es necessita que entris les dades. Crea un arxiu anomenat `data.tsv`.
Cada fila ha de contindre les següents dades separades per tabuladors:
```
NOM	PROVEIDOR	UNITATS		ARGUMENTS...
```
Aquest és el significat de cada valor:
- Nom és el nom del producte, dona igual el que posis.
- Proveidor és el supermercat d'on s'obté. Mira la llista de proveidors més abaix.
- Unitats és el número d'unitats. Si vols saber el preu per pastanaga, i cada pack en té 10, posa 10 aquí.
- Arguments és informació que el proveidor necessita. Poden ser múltiples dades separades per tabuladors.

## Proveidors

- Bonpreu: Obté la informació de la web. Per exemple una poma té la URL https://www.compraonline.bonpreuesclat.cat/products/90041/details.
    - Arguments: El codi de la URL. En el cas de la poma es tracta de `90041`
- Mercadona: Obté la informació de la web. Per exemple una poma té la URL https://tienda.mercadona.es/product/8177/manzana-roja-dulce-pieza.
    - Arguments: CODI	UBICACIÓ
    - El codi és `8177` en el cas de la poma.
    - La ubicació és més difícil de trobar. 
       - Per fer-ho fàcil pots utilitzar `bcn1`.
       - Si vols fer-ho complicat, entra a la Web de compra Online de Mercadona i no entris el teu codi postal. Inspecciona la web, ves a la pestanya `Network` i aleshores entra el teu codi postal. Tornant a `Network`, si filtres per `home/?` trobaràs una `Payload` amb `lang` i `wh`. El codi de la teva ubicació és l'indicat a `wh`.

### Exemple
Tens un exemple [aquí](./end-to-end/example.tsv).
```tsv
Fruita	Bonpreu	1	90041
Alvocat	Bonpreu	1	82055
Fullaraca	Bonpreu	1	61688
Pastanaga	Bonpreu	10	15297
Pebrot verd	Bonpreu	1	46754
...
```

## Executció
Un cop compilat pots executar-ho:
```
./bin/compra -i data.tsv
```
Trigarà uns segons en executar però produirà un resultat com el següent:
```tsv
Product               Price
Fruita                0,65 €
Alvocat               1,39 €
Fullaraca             1,15 €
Pastanaga             0,17 €
Pebrot verd           0,49 €
...
```

Si vols especificar el format o altres opcions, llegeix l'ajuda
```
./bin/compra -help
```

