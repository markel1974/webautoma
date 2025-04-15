package help

import (
	"fmt"
	"sort"
	"strings"
)

// Provider fornisce accesso diretto alle stringhe di aiuto predefinite per i comandi.
type Provider struct {
	// Mappa che contiene il testo di aiuto per ogni comando.
	// Chiave: nome del comando (es. "open", "selectWindow")
	// Valore: stringa multi-linea con la descrizione, parametri ed esempio.
	commandHelp map[string]string
}

// NewHelp crea e inizializza un nuovo Provider.
// Le stringhe di aiuto sono definite direttamente qui.
func NewHelp() *Provider {
	helpMap := make(map[string]string)

	// Popola la mappa con le descrizioni che abbiamo generato
	// (Uso backtick ` per stringhe multi-linea)

	// --- Categoria: Navigazione ---
	helpMap["open"] = `
        - Descrizione: Apre un URL nel browser o naviga verso un percorso relativo all'URL base corrente.
        - Parametri:
          - target: L'URL completo (es. https://google.com) o un percorso relativo (es. /pagina). Se relativo, viene aggiunto all'url base definito nel file .side o all'ultimo URL base noto.
          - value: Non utilizzato.
        - Esempio:
          {
            "id": "...",
            "command": "open",
            "target": "https://www.google.com",
            "value": ""
          }
    `

	// --- Categoria: Interazione con Elementi ---
	helpMap["click"] = `
        - Descrizione: Simula un click del mouse sull'elemento specificato. Attende che l'elemento sia pronto (visibile e abilitato) prima di cliccare.
        - Parametri:
          - target: Il selettore dell'elemento da cliccare (es. id=myButton, css=.submit-btn).
          - value: Non utilizzato.
          - until: Se impostato a 1 o 2, attende che l'elemento sia pronto prima del click.
        - Esempio:
          {
            "id": "...",
            "command": "click",
            "target": "id=loginButton",
            "value": "",
            "until": "1"
          }
    `
	helpMap["type"] = `
       - Descrizione: Inserisce del testo in un campo input o textarea. Simula la digitazione carattere per carattere con un piccolo ritardo (vedi humanWait).
       - Parametri:
         - target: Il selettore dell'elemento in cui scrivere (es. id=username, name=password).
         - value: Il testo da inserire. Supporta variabili (es. {{.mioUsername}}).
         - until: Se impostato a 1 o 2, attende che l'elemento sia pronto prima di scrivere.
       - Esempio:
         {
           "id": "...",
           "command": "type",
           "target": "id=searchField",
           "value": "Testo da cercare"
         }
         {
           "id": "...",
           "command": "type",
           "target": "name=password",
           "value": "{{.userPassword}}"
         }
    `
	helpMap["select"] = `
       - Descrizione: Seleziona un'opzione da un elemento <select> (dropdown) basandosi sul testo visibile dell'opzione (label).
       - Parametri:
         - target: Il selettore dell'elemento <select>.
         - value: La stringa label=Testo Dell'Opzione che identifica l'opzione da selezionare.
         - until: Se impostato a 1 o 2, attende che l'elemento <select> sia pronto.
       - Esempio:
         {
           "id": "...",
           "command": "select",
           "target": "id=countryDropdown",
           "value": "label=Italia"
         }
    `

	// --- Categoria: Gestione Timer (Custom webautoma) ---
	helpMap["timerCreate"] = `
       - Descrizione: Crea e inizializza un nuovo timer, senza farlo partire. Utile per misurare tempi composti da più azioni.
       - Parametri:
         - target: L'ID univoco da assegnare al timer (es. loginTime).
         - value: Una descrizione opzionale per il timer (riportata nei log).
       - Esempio:
         {
           "id": "...",
           "command": "timerCreate",
           "target": "pageLoadTimer",
           "value": "Tempo caricamento pagina iniziale"
         }
    `
	helpMap["timerStart"] = `
       - Descrizione: Avvia (o riavvia) un timer precedentemente creato con timerCreate. Registra il tempo di inizio.
       - Parametri:
         - target: L'ID del timer da avviare.
         - value: Non utilizzato.
       - Esempio:
         {
           "id": "...",
           "command": "timerStart",
           "target": "pageLoadTimer",
           "value": ""
         }
    `
	helpMap["timerStop"] = `
       - Descrizione: Ferma un timer precedentemente avviato. Registra l'intervallo trascorso dall'ultimo timerStart o timerStop. Se value è "finalize", finalizza il timer e scrive l'evento completo nel log JSON; altrimenti, registra solo l'intervallo parziale.
       - Parametri:
         - target: L'ID del timer da fermare.
         - value: Se impostato a finalize (case-insensitive), finalizza il timer. Altrimenti, non fa nulla di speciale oltre a fermare l'intervallo corrente.
       - Esempio (Stop parziale):
         {
           "id": "...",
           "command": "timerStop",
           "target": "userActionTimer",
           "value": ""
         }
       - Esempio (Stop e Finalize):
         {
           "id": "...",
           "command": "timerStop",
           "target": "totalTestTimer",
           "value": "finalize"
         }
    `
	helpMap["timerFinalize"] = `
       - Descrizione: Finalizza un timer, calcolando il tempo totale trascorso sommando tutti gli intervalli registrati con timerStop. Scrive l'evento completo nel log JSON. Il timer non può più essere usato dopo la finalizzazione.
       - Parametri:
         - target: L'ID del timer da finalizzare.
         - value: Non utilizzato.
       - Esempio:
         {
           "id": "...",
           "command": "timerFinalize",
           "target": "loginProcessTimer",
           "value": ""
         }
    `

	// --- Categoria: Gestione Stack Variabili (Custom webautoma) ---
	helpMap["stackAdd"] = `
       - Descrizione: Trova un elemento, ne estrae alcune proprietà (testo, tag, stato visualizzato/abilitato) e le salva in una mappa interna ("stack") associandole all'ID fornito nel campo id del comando. Utile per memorizzare stati intermedi o valori dinamici.
       - Parametri:
         - target: Il selettore dell'elemento da cui estrarre le informazioni.
         - value: Non utilizzato direttamente.
         - id: (Campo standard del comando) Importante: Questo id viene usato come chiave per memorizzare le informazioni nello stack.
         - until: Se impostato a 1, attende che l'elemento sia pronto.
       - Esempio:
         {
           "id": "userInfo", // Questo ID sarà la chiave nello stack
           "command": "stackAdd",
           "target": "id=userDetails",
           "value": "",
           "until": "1"
         }
    `
	helpMap["stackReset"] = `
        - Descrizione: Cancella completamente lo stack interno delle variabili.
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "stackReset",
            "target": "",
            "value": ""
          }
    `
	helpMap["stackPrint"] = `
        - Descrizione: Stampa il contenuto corrente dello stack sulla console (output standard) in formato JSON indentato. Utile per debugging.
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "stackPrint",
            "target": "",
            "value": ""
          }
    `

	// --- Categoria: Controllo Flusso & Utility ---
	helpMap["jump"] = `
        - Descrizione: Salta l'esecuzione a un altro comando all'interno dello stesso test, identificato dal suo id. Nota: Questo non è un comando standard Selenium IDE e rende il flusso del test più difficile da seguire nell'IDE stesso.
        - Parametri:
          - target: L'id del comando a cui saltare.
          - value: Non utilizzato.
        - Esempio:
          {
            "id": "...",
            "command": "jump",
            "target": "inizioCiclo", // Salta al comando con id "inizioCiclo"
            "value": ""
          }
    `
	helpMap["pause"] = `
        - Descrizione: Sospende l'esecuzione per un numero specificato di millisecondi.
        - Parametri:
          - target: Il numero di millisecondi per cui sospendere l'esecuzione.
          - value: Non utilizzato.
        - Esempio:
          {
            "id": "...",
            "command": "pause",
            "target": "5000", // Pausa per 5 secondi
            "value": ""
          }
    `
	helpMap["humanWait"] = `
        - Descrizione: Introduce una pausa "umana", ovvero una pausa di durata variabile casuale basata su un valore base configurabile (parametro -humanWait o executor.humanWaitBase nel codice). Se viene fornito un valore nel value del comando, usa quello come durata fissa in millisecondi.
        - Parametri:
          - target: Non utilizzato.
          - value: Durata fissa della pausa in millisecondi (opzionale). Se omesso, usa la pausa casuale basata su humanWaitBase.
        - Esempio (Pausa casuale):
          {
            "id": "...",
            "command": "humanWait",
            "target": "",
            "value": ""
          }
        - Esempio (Pausa fissa):
          {
            "id": "...",
            "command": "humanWait",
            "target": "",
            "value": "1500" // Pausa fissa di 1.5 secondi
          }
    `

	// --- Categoria: Asserzioni e Verifiche ---
	helpMap["assert"] = `
        - Descrizione: Verifica che il testo visibile di un elemento corrisponda esattamente al valore fornito. Fallisce se il testo è diverso.
        - Parametri:
          - target: Il selettore dell'elemento il cui testo deve essere verificato (es. id=messaggioErrore, css=.risultato).
          - value: Il testo esatto che ci si aspetta di trovare nell'elemento. Supporta variabili {{.variabile}}.
          - until: (Opzionale) Se impostato a 1 o 2, attende che l'elemento sia pronto (visibile e abilitato) prima di effettuare la verifica.
        - Esempio:
          {
            "id": "...",
            "command": "assert",
            "target": "id=statusMessage",
            "value": "Operazione completata.",
            "until": "1"
          }
    `
	helpMap["exists"] = `
        - Descrizione: Verifica che un elemento specificato esista nel DOM e sia pronto (visibile e abilitato) entro il timeout configurato. Fallisce se l'elemento non viene trovato o non diventa pronto.
        - Parametri:
          - target: Il selettore dell'elemento da cercare (es. id=confermaPopup, css=button.primary).
          - value: Non utilizzato.
          - until: (Opzionale) Se impostato a 1 o 2, attende attivamente che l'elemento esista e sia pronto per la durata del timeout. Se omesso (o 0), verifica solo se l'elemento è presente e pronto nello stato attuale della pagina (meno comune per verifiche robuste).
        - Esempio:
          {
            "id": "...",
            "command": "exists",
            "target": "css=.loading-spinner",
            "value": "",
            "until": "1" // Attende che lo spinner sia visibile/abilitato
          }
    `
	helpMap["until"] = `
        - Descrizione: Attende fino a quando un elemento specificato *non* è più presente o *non* è più pronto (visibile/abilitato) sulla pagina, oppure fino allo scadere del timeout. Utile per aspettare la scomparsa di elementi temporanei (es. messaggi di caricamento). Fallisce se l'elemento rimane presente e pronto alla fine del timeout.
        - Parametri:
          - target: Il selettore dell'elemento da monitorare per la sua scomparsa o inattività.
          - value: Non utilizzato.
          - until: (Opzionale, ma **consigliato impostarlo a 1** per questo comando) Se 1, attende che l'elemento *non* sia pronto/presente. Se 0 o 2, attende che sia pronto, ma il comando fallirà se l'elemento *è* effettivamente pronto (uso meno intuitivo).
        - Esempio:
          {
            "id": "...",
            "command": "until",
            "target": "css=.loading-indicator",
            "value": "",
            "until": "1" // Attende che l'indicatore di caricamento scompaia o diventi non pronto
          }
    `

	// --- Categoria: Gestione Finestre/Frame/Alert ---
	helpMap["setWindowSize"] = `
        - Descrizione: Ridimensiona la finestra corrente alle dimensioni specificate.
        - Parametri:
          - target: La dimensione desiderata nel formato "LarghezzaxAltezza" (es. "1280x720").
          - value: Non utilizzato.
        - Esempio:
          {
            "id": "...",
            "command": "setWindowSize",
            "target": "1920x1080",
            "value": ""
          }
    `
	helpMap["selectWindow"] = `
        - Descrizione: Sposta il focus del driver su una finestra o tab specifico. Può usare l'handle diretto del WebDriver, un nome assegnato precedentemente con storeWindowHandle (usando la sintassi ${nomeHandle}), o un indice numerico (meno comune/affidabile).
        - Parametri:
          - target: L'identificatore della finestra. Formati possibili:
              - handle=NOME_HANDLE_WEBDRIVER (raramente usato manualmente)
              - ${nomeHandleSalvato} (nome assegnato con storeWindowHandle)
              - Potenzialmente un indice numerico (da verificare nel codice WebDriver)
          - value: Non utilizzato.
        - Esempio (usando un handle salvato):
          {
            "id": "...",
            "command": "selectWindow",
            "target": "${finestraPrincipale}",
            "value": ""
          }
    `
	helpMap["selectWindowMain"] = `
        - Descrizione: Riporta il focus sulla finestra principale/iniziale (quella aperta all'avvio o impostata con setWindowMain).
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "selectWindowMain",
            "target": "",
            "value": ""
          }
    `
	helpMap["selectWindowTitle"] = `
        - Descrizione: Sposta il focus sulla prima finestra/tab il cui titolo corrisponde (parzialmente, esattamente, o tramite regex) al value specificato.
        - Parametri:
          - target: Modalità di confronto del titolo (Opzionale, default: contains):
              - contains: Il titolo della finestra contiene value (case-insensitive).
              - exact: Il titolo della finestra è esattamente value.
              - regexp: Il titolo della finestra matcha l'espressione regolare in value.
          - value: Il testo del titolo (o l'espressione regolare) da cercare.
        - Esempio (Contains):
          {
            "id": "...",
            "command": "selectWindowTitle",
            "target": "contains", // o omesso
            "value": "Pagina Risultati"
          }
        - Esempio (Regexp):
          {
            "id": "...",
            "command": "selectWindowTitle",
            "target": "regexp",
            "value": "^Carrello \\(\\d+\\)$" // Es: Titolo "Carrello (3)"
          }
    `
	helpMap["closeWindow"] = `
        - Descrizione: Chiude la finestra o il tab *attualmente* in focus. Non può chiudere la finestra principale/iniziale. Dopo la chiusura, il focus *non* viene spostato automaticamente; usare selectWindowMain o un altro comando selectWindow* se necessario.
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "closeWindow",
            "target": "",
            "value": ""
          }
    `
	helpMap["storeWindowHandle"] = `
        - Descrizione: Salva l'handle della finestra attualmente in focus associandolo a un nome simbolico. Questo nome può essere usato successivamente in selectWindow o close con la sintassi ${nomeHandle}.
        - Parametri:
          - target: Il nome da assegnare all'handle della finestra corrente (es. finestraLogin, popupDettaglio).
          - value: Non utilizzato.
        - Campi Specifici: Può utilizzare windowHandleName (ridondante?), windowTimeout, opensWindow per gestire l'attesa dell'apertura di una nuova finestra prima di salvarne l'handle (se opensWindow è true).
        - Esempio (Salvataggio handle corrente):
          {
            "id": "...",
            "command": "storeWindowHandle",
            "target": "mainWindow", // Salva l'handle corrente come "mainWindow"
            "value": ""
          }
        - Esempio (Attesa e salvataggio popup):
          {
            "id": "...",
            "command": "storeWindowHandle",
            "target": "finestraPopup",
            "value": "",
            "opensWindow": true,
            "windowTimeout": 5000 // Attende fino a 5s che appaia una nuova finestra
          }
    `
	helpMap["windowHandles"] = `
        - Descrizione: Stampa sulla console (output standard) l'elenco degli handle di tutte le finestre attualmente aperte. Utile principalmente per debugging.
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "windowHandles",
            "target": "",
            "value": ""
          }
    `
	helpMap["close"] = `
        - Descrizione: Chiude una specifica finestra identificata dal suo handle o da un nome precedentemente salvato con storeWindowHandle. A differenza di closeWindow, questo comando richiede l'identificatore della finestra da chiudere.
        - Parametri:
          - target: L'identificatore della finestra da chiudere (es. ${popup}, handle=HANDLE_ID).
          - value: Non utilizzato.
        - Esempio:
          {
            "id": "...",
            "command": "close",
            "target": "${finestraDaChiudere}",
            "value": ""
          }
    `
	helpMap["setWindowMain"] = `
        - Descrizione: Imposta l'handle della finestra *attualmente* in focus come la nuova finestra "principale" o "root" per webautoma. Utile se il flusso di lavoro si sposta permanentemente su una nuova finestra principale.
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "setWindowMain",
            "target": "",
            "value": ""
          }
    `
	helpMap["selectFrame"] = `
        - Descrizione: Sposta il focus su un frame (o iframe) all'interno della pagina corrente.
        - Parametri:
          - target: Identificatore del frame:
              - index=N: Indice numerico del frame (base 0).
              - relative=parent: Passa al frame genitore.
              - relative=top: Passa al contesto principale della pagina (fuori da tutti i frame).
              - Stringa: Nome o ID dell'elemento (i)frame.
              - Vuoto: Resetta al contesto principale della pagina (equivalente a relative=top).
          - value: Non utilizzato.
        - Esempio (Per indice):
          {
            "id": "...",
            "command": "selectFrame",
            "target": "index=0",
            "value": ""
          }
        - Esempio (Per ID/Nome):
          {
            "id": "...",
            "command": "selectFrame",
            "target": "contentFrame",
            "value": ""
          }
    `
	helpMap["selectParentFrame"] = `
        - Descrizione: Sposta il focus dal frame corrente al suo frame genitore diretto. Equivalente a selectFrame con target=relative=parent.
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "selectParentFrame",
            "target": "",
            "value": ""
          }
    `
	helpMap["selectAlert"] = `
        - Descrizione: Gestisce un alert JavaScript (popup alert(), confirm(), prompt()) che appare sulla pagina. Legge il testo dell'alert (per log) e poi lo accetta o lo chiude.
        - Parametri:
          - target: Non utilizzato.
          - value: Determina l'azione:
              - 1: Accetta l'alert (es. preme "OK" in confirm).
              - 0 o omesso: Chiude/Annulla l'alert (es. preme "Annulla" in confirm).
        - Esempio (Accetta):
          {
            "id": "...",
            "command": "selectAlert",
            "target": "",
            "value": "1"
          }
        - Esempio (Chiude):
          {
            "id": "...",
            "command": "selectAlert",
            "target": "",
            "value": "0" // o omesso
          }
    `

	// --- Categoria: Interazioni Avanzate (Mouse) ---
	helpMap["rightClick"] = `
        - Descrizione: Simula un click con il tasto destro del mouse sull'elemento specificato. Attende che l'elemento sia pronto.
        - Parametri:
          - target: Il selettore dell'elemento su cui fare click destro.
          - value: Non utilizzato.
          - until: (Opzionale) Se 1 o 2, attende che l'elemento sia pronto.
        - Esempio:
          {
            "id": "...",
            "command": "rightClick",
            "target": "id=contextMenuArea",
            "value": "",
            "until": "1"
          }
    `
	helpMap["doubleClick"] = `
        - Descrizione: Simula un doppio click con il tasto sinistro del mouse sull'elemento specificato. Attende che l'elemento sia pronto.
        - Parametri:
          - target: Il selettore dell'elemento su cui fare doppio click.
          - value: Non utilizzato.
          - until: (Opzionale) Se 1 o 2, attende che l'elemento sia pronto.
        - Esempio:
          {
            "id": "...",
            "command": "doubleClick",
            "target": "css=div.editable",
            "value": "",
            "until": "1"
          }
    `
	helpMap["mouseOver"] = `
        - Descrizione: Sposta il cursore del mouse sopra l'elemento specificato, potenzialmente attivando effetti hover o tooltip. Attende che l'elemento sia pronto.
        - Parametri:
          - target: Il selettore dell'elemento su cui spostare il mouse.
          - value: Non utilizzato.
          - until: (Opzionale) Se 1 o 2, attende che l'elemento sia pronto.
        - Esempio:
          {
            "id": "...",
            "command": "mouseOver",
            "target": "id=menuItem",
            "value": "",
            "until": "1"
          }
    `
	helpMap["mouseOut"] = `
        - Descrizione: Simula l'evento del mouse che esce dall'area dell'elemento specificato. Nota: Nel codice attuale (doMouse), questo comando non sembra eseguire azioni specifiche sul WebDriver.
        - Parametri:
          - target: Il selettore dell'elemento da cui il mouse "esce".
          - value: Non utilizzato.
        - Esempio:
          {
            "id": "...",
            "command": "mouseOut",
            "target": "id=menuItem",
            "value": ""
          }
    `
	helpMap["mouseDownAt"] = `
        - Descrizione: Simula la pressione (senza rilascio) del tasto sinistro del mouse sull'elemento specificato, potenzialmente a coordinate relative. Inizia un'operazione di trascinamento (drag).
        - Parametri:
          - target: Il selettore dell'elemento su cui premere il mouse.
          - value: (Opzionale) Coordinate relative all'angolo in alto a sinistra dell'elemento, formato "X,Y" (es. "10,15"). Se omesso, preme al centro (o default del driver).
          - until: (Opzionale) Se 1 o 2, attende che l'elemento sia pronto.
        - Esempio:
          {
            "id": "...",
            "command": "mouseDownAt",
            "target": "id=draggableElement",
            "value": "5,5", // Premi vicino all'angolo in alto a sinistra
            "until": "1"
          }
    `
	helpMap["mouseMoveAt"] = `
        - Descrizione: Sposta il mouse mentre il tasto sinistro è tenuto premuto (iniziato con mouseDownAt). Usato per il trascinamento (drag). Richiede un mouseDownAt precedente sullo stesso elemento implicito.
        - Parametri:
          - target: (Generalmente ignorato, agisce sull'elemento del mouseDownAt).
          - value: Coordinate relative all'angolo in alto a sinistra dell'elemento *originale* del mouseDownAt, formato "X,Y". Indica la posizione *a cui* spostare il mouse.
        - Esempio:
          {
            "id": "...",
            "command": "mouseMoveAt",
            "target": "", // Target non necessario qui
            "value": "100,50" // Sposta il mouse a 100px destra, 50px sotto dall'origine del drag
          }
    `
	helpMap["mouseMultipleMoveAt"] = `
        - Descrizione: Simile a mouseMoveAt, ma esegue una sequenza di spostamenti relativi consecutivi mentre il tasto è premuto. Utile per simulare un trascinamento lungo un percorso. Richiede un mouseDownAt precedente.
        - Parametri:
          - target: (Generalmente ignorato).
          - value: Sequenza di coordinate relative *incrementali*, separate da |. Ogni coppia "X,Y" è relativa alla *posizione precedente*. Formato: "dX1,dY1|dX2,dY2|...".
        - Esempio:
          {
            "id": "...",
            "command": "mouseMultipleMoveAt",
            "target": "",
            "value": "50,0|0,50|-50,0" // Sposta 50px a destra, poi 50px in basso, poi 50px a sinistra
          }
    `
	helpMap["mouseUpAt"] = `
        - Descrizione: Simula il rilascio del tasto sinistro del mouse, completando un'operazione di trascinamento (drag and drop). Richiede un mouseDownAt precedente.
        - Parametri:
          - target: (Generalmente ignorato, agisce sull'elemento del mouseDownAt).
          - value: (Opzionale) Coordinate relative all'angolo in alto a sinistra dell'elemento *originale* del mouseDownAt, formato "X,Y". Indica la posizione *finale* in cui rilasciare il mouse. Se omesso, rilascia nella posizione corrente.
        - Esempio:
          {
            "id": "...",
            "command": "mouseUpAt",
            "target": "",
            "value": "200,100" // Rilascia il mouse a 200px destra, 100px sotto dall'origine del drag
          }
    `

	// --- Categoria: Interazioni Avanzate (Tastiera - Custom webautoma) ---
	helpMap["actionsSendKeys"] = `
        - Descrizione: Invia la pressione di un tasto speciale (non alfanumerico) o una sequenza semplice all'elemento attivo nella pagina. Utile per simulare Invio, Tab, Frecce direzionali, ecc.
        - Parametri:
          - target: (Generalmente non usato, agisce sull'elemento attivo).
          - value: Il nome del tasto speciale (es. Enter, Tab, ArrowDown, Control, Alt, Shift, F5, etc. - vedi wd/base/keys.go per la lista completa dei nomi mappati) oppure una sequenza di caratteri semplici. Il mapping cerca KeyFromMapping in base/keys.go per i nomi speciali.
        - Esempio (Invio):
          {
            "id": "...",
            "command": "actionsSendKeys",
            "target": "",
            "value": "Enter"
          }
        - Esempio (Freccia Giù):
          {
            "id": "...",
            "command": "actionsSendKeys",
            "target": "",
            "value": "ArrowDown"
          }
    `
	helpMap["keys"] = `
        - Descrizione: Permette di definire sequenze complesse di interazioni da tastiera, includendo pressione (press), rilascio (release) e pause (pause). Utile per simulare combinazioni di tasti (es. Ctrl+C) o comportamenti specifici.
        - Parametri:
          - target: (Opzionale) Selettore dell'elemento a cui inviare gli eventi. Se omesso, invia all'elemento attivo.
          - value: Stringa formattata che descrive la sequenza, separata da virgola. Ogni elemento è azione:valore.
              - press:KEY: Simula la pressione di un tasto (es. press:Control, press:c). Usa KeyFromMapping per i tasti speciali.
              - release:KEY: Simula il rilascio di un tasto (es. release:Control, release:c).
              - pause:MS: Inserisce una pausa in millisecondi (es. pause:100).
        - Esempio (Ctrl+A, Ctrl+C):
          {
            "id": "...",
            "command": "keys",
            "target": "id=myTextArea",
            "value": "press:Control,press:a,release:a,pause:50,press:c,release:c,release:Control"
          }
        - Esempio (Scrivere "test" tenendo premuto Shift):
          {
            "id": "...",
            "command": "keys",
            "target": "id=myInput",
            "value": "press:Shift,press:t,release:t,press:e,release:e,press:s,release:s,press:t,release:t,release:Shift"
          }
    `

	// --- Categoria: File & Download (Custom webautoma) ---
	helpMap["download"] = `
        - Descrizione: Scarica un file direttamente da un URL specificato e lo salva nel percorso locale indicato. Utilizza i cookie di sessione correnti del browser per gestire eventuali autenticazioni richieste dal server per il download.
        - Parametri:
          - target: L'URL completo da cui scaricare il file. Se l'URL inizia con '/', viene considerato relativo all'URL base della pagina corrente (ottenuto dall'ultimo comando 'open' o navigazione). Supporta variabili {{.variabile}}.
          - value: Il percorso completo, incluso il nome del file, dove salvare il file scaricato sul sistema locale (es. /percorso/locale/nomefile.pdf o C:\Download\report.xlsx). Supporta variabili {{.variabile}}.
        - Esempio:
          {
            "id": "...",
            "command": "download",
            "target": "https://example.com/resources/documento.zip",
            "value": "/home/utente/download/archivio.zip"
          }
          {
            "id": "...",
            "command": "download",
            "target": "/api/export?id={{.reportId}}", // Relativo all'URL base
            "value": "report_{{.reportId}}.csv"
          }
    `
	helpMap["clickDownload"] = `
        - Descrizione: Trova un elemento sulla pagina (solitamente un link <a>), ne estrae l'URL dall'attributo href, quindi scarica il file da quell'URL nel percorso locale specificato. Utile quando l'URL del download è dinamico o non noto a priori. Utilizza i cookie di sessione correnti.
        - Parametri:
          - target: Il selettore dell'elemento (es. link <a>) che contiene l'attributo href con l'URL del file da scaricare.
          - value: Il percorso completo, incluso il nome del file, dove salvare il file scaricato sul sistema locale. Supporta variabili {{.variabile}}.
          - until: (Opzionale) Se impostato a 1 o 2, attende che l'elemento sia pronto prima di tentare di leggerne l'attributo href.
        - Esempio:
          {
            "id": "...",
            "command": "clickDownload",
            "target": "css=a.download-link[data-file='report']",
            "value": "/tmp/report_scaricato.pdf",
            "until": "1"
          }
    `

	// --- Categoria: Scrolling (Custom webautoma) ---
	helpMap["scroll"] = `
        - Descrizione: Esegue uno scroll della pagina (viewport) o a partire da un elemento specifico, di una determinata quantità (delta X, delta Y). Utile per spostare la visuale di una quantità fissa. Può simulare uno scroll fluido suddividendolo in passi con pause intermedie. Utilizza l'Actions API del WebDriver (simulazione rotellina/touchpad).
        - Parametri:
          - target: (Opzionale) Il selettore dell'elemento da cui calcolare il punto di partenza dello scroll. Se omesso, lo scroll parte dalle coordinate specificate nei campi x/y del comando (default 0,0) più eventuali offsetX/offsetY.
          - value: Specifica lo spostamento (delta) e opzionalmente i passi e l'intervallo. Formato: "DeltaX,DeltaY[,NumeroPassi,IntervalloMs]".
              - DeltaX, DeltaY: Pixel di spostamento orizzontale e verticale (negativi per sinistra/su, positivi per destra/giù).
              - NumeroPassi: (Opzionale) Numero di passi in cui suddividere lo scroll totale.
              - IntervalloMs: (Opzionale, richiede NumeroPassi) Millisecondi di pausa tra un passo e l'altro.
          - x, y: (Opzionali, usati solo se target è omesso) Coordinate X, Y *assolute* nella viewport da cui inizia l'azione di scroll. Default: 0,0.
          - offsetX, offsetY: (Opzionali) Offset in pixel da aggiungere alle coordinate di partenza (sia quelle x/y sia quelle calcolate dall'elemento target).
        - Esempio (Scroll viewport di 500px in basso):
          {
            "id": "...",
            "command": "scroll",
            "target": "", // Scroll della viewport
            "value": "0,500"
          }
        - Esempio (Scroll viewport fluido di 1000px in basso in 10 passi):
          {
            "id": "...",
            "command": "scroll",
            "target": "",
            "value": "0,1000,10,50" // 10 passi, 50ms di pausa tra passi
          }
        - Esempio (Scroll partendo dal basso di un elemento header):
          {
            "id": "...",
            "command": "scroll",
            "target": "id=mainHeader",
            "value": "0,300", // Scrolla 300px in basso
            "offsetY": 50 // Partendo 50px sotto l'header
          }
    `
	helpMap["scrollTo"] = `
        - Descrizione: Scrolla la pagina in modo che un punto specifico *all'interno* di un elemento target diventi visibile nella viewport. Se l'elemento non è inizialmente visibile, prova prima a portarlo in vista. Può simulare uno scroll fluido. Utilizza l'Actions API del WebDriver.
        - Parametri:
          - target: (Opzionale) Il selettore dell'elemento target verso cui scrollare. Se omesso, utilizza l'elemento attualmente attivo (quello con il focus).
          - value: Specifica le coordinate *relative all'angolo in alto a sinistra dell'elemento target* e opzionalmente i passi e l'intervallo. Formato: "TargetX,TargetY[,NumeroPassi,IntervalloMs]".
              - TargetX, TargetY: Coordinate X, Y *dentro* l'elemento target che si desidera portare in vista. "0,0" corrisponde all'angolo in alto a sinistra dell'elemento.
              - NumeroPassi: (Opzionale) Numero di passi per raggiungere la posizione.
              - IntervalloMs: (Opzionale, richiede NumeroPassi) Millisecondi di pausa tra i passi.
          - offsetX, offsetY: (Opzionali) Offset in pixel da aggiungere alle coordinate TargetX, TargetY specificate in value.
          - until: (Opzionale) Se impostato a 1 o 2, attende che l'elemento target sia pronto prima di tentare lo scroll.
        - Esempio (Scrollare fino all'inizio di un footer):
          {
            "id": "...",
            "command": "scrollTo",
            "target": "id=pageFooter",
            "value": "0,0", // Porta l'angolo 0,0 del footer in vista
            "until": "1"
          }
        - Esempio (Scrollare fluidamente a un punto specifico dentro un div):
          {
            "id": "...",
            "command": "scrollTo",
            "target": "css=div.scrollable-content",
            "value": "0,500,10,50" // Porta il punto Y=500 dentro il div in vista, in 10 passi
          }
    `

	// --- Categoria: Custom Avanzati (webautoma) ---
	helpMap["otp"] = `
        - Descrizione: Recupera un One-Time Password (OTP) da un account email. Si connette al server IMAP specificato, cerca l'email più recente che soddisfa i criteri (oggetto, età massima), estrae il codice OTP dal corpo dell'email tramite un'espressione regolare, e lo memorizza internamente. L'OTP recuperato può essere utilizzato nei comandi successivi (ad esempio 'type') usando la variabile {{.otp}}.
        - Parametri:
          - target: La stringa di connessione al server IMAP. Formato: protocollo[authMode]://utente:password@server:porta
              - protocollo: Può essere tls (consigliato), starttls, o insecure.
              - [authMode]: (Opzionale) Specificare [oauth] o [oauth2] se si utilizza l'autenticazione OAuth/OAuth2 invece della password diretta.
              - utente: Nome utente per l'accesso IMAP.
              - password: Password o token OAuth.
              - server: Indirizzo del server IMAP.
              - porta: Porta del server IMAP (es. 993 per TLS, 143 per StartTLS/Insecure).
              - Esempio TLS: tls://mia.email@example.com:LaMiaPassword@imap.example.com:993
              - Esempio OAuth2: tls[oauth2]://utente@gmail.com:TokenDiAccessoOAuth2@imap.gmail.com:993
          - value: Stringa di configurazione per la ricerca e l'estrazione, con parametri separati da |||. Formato: "RegExpOggetto|||RegExpCorpoConGruppoDiCatturaOTP[|||IntervalloVerificaSec[|||ValiditaEmailMin]]"
              - RegExpOggetto: Espressione regolare (Go standard) per identificare l'oggetto dell'email contenente l'OTP (es. ^Codice di verifica.*$).
              - RegExpCorpoConGruppoDiCatturaOTP: Espressione regolare (Go standard) applicata al corpo dell'email per estrarre l'OTP. **Deve** contenere un gruppo di cattura tra parentesi () che isoli esattamente il codice OTP (es. Il tuo codice OTP è ([0-9]{6})\., cattura 6 cifre).
              - IntervalloVerificaSec: (Opzionale, default: 60) Numero massimo di secondi durante i quali webautoma tenterà di recuperare l'email (controllando periodicamente la casella).
              - ValiditaEmailMin: (Opzionale, default: 5) Età massima in minuti che l'email può avere per essere considerata valida.
        - Risultato: L'OTP estratto viene salvato nella variabile interna otp, accessibile come {{.otp}} nei campi value dei comandi successivi.
        - Esempio:
          {
            "id": "recuperaOTP",
            "command": "otp",
            "target": "tls://utente@example.com:passwordSegreta@imap.example.com:993",
            "value": "Il tuo codice monouso Esempio Corp|||codice di verifica: ([A-Z0-9]+)|||90|||3"
            // Cerca email con oggetto "Il tuo codice monouso Esempio Corp"
            // Estrae un codice alfanumerico dal corpo (es. "codice di verifica: XY78Z1")
            // Tenta per 90 secondi, email valide se più recenti di 3 minuti
          }
          {
            "id": "inserisciOTP",
            "command": "type",
            "target": "id=otpField",
            "value": "{{.otp}}" // Usa l'OTP recuperato nel passo precedente
          }
    `

	// --- Categoria: Debug e Meta-Comandi (Custom webautoma) ---
	helpMap["disableError"] = `
        - Descrizione: Disabilita temporaneamente l'interruzione dell'esecuzione in caso di errore nei comandi successivi. Utile per tentare azioni che potrebbero fallire senza bloccare l'intero test. L'errore verrà comunque loggato (a meno che non sia disabilitato anche il debug). Usare enableError per riattivare il comportamento normale.
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "disableError",
            "target": "",
            "value": ""
          }
    `
	helpMap["enableError"] = `
        - Descrizione: Riattiva l'interruzione dell'esecuzione in caso di errore, annullando l'effetto di un precedente disableError. Questo è il comportamento predefinito.
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "enableError",
            "target": "",
            "value": ""
          }
    `
	helpMap["disableDebug"] = `
        - Descrizione: Disabilita l'output di log dettagliato (livello debug) generato dall'adapter webautoma durante l'esecuzione.
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "disableDebug",
            "target": "",
            "value": ""
          }
    `
	helpMap["enableDebug"] = `
        - Descrizione: Riattiva l'output di log dettagliato (livello debug). Utile se è stato precedentemente disabilitato.
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "enableDebug",
            "target": "",
            "value": ""
          }
    `
	helpMap["status"] = `
        - Descrizione: Richiede e stampa sulla console (output standard) le informazioni di stato del server WebDriver (versione, OS, disponibilità). Utile per diagnosi.
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "status",
            "target": "",
            "value": ""
          }
    `
	helpMap["execId"] = `
        - Descrizione: Imposta un identificatore globale per l'intera esecuzione corrente. Questo ID verrà incluso nei log JSON generati, utile per correlare eventi di diverse esecuzioni.
        - Parametri:
          - target: La stringa da usare come ID dell'esecuzione.
          - value: Non utilizzato.
        - Esempio:
          {
            "id": "...",
            "command": "execId",
            "target": "run_produzione_sera",
            "value": ""
          }
    `
	helpMap["probe"] = `
        - Descrizione: Imposta un identificatore per la "sonda" o istanza specifica del test in corso. Questo ID verrà incluso nei log JSON, utile per distinguere i risultati quando più istanze dello stesso test girano in parallelo o per identificare specifici punti di monitoraggio.
        - Parametri:
          - target: La stringa da usare come ID della sonda/istanza.
          - value: Non utilizzato.
        - Esempio:
          {
            "id": "...",
            "command": "probe",
            "target": "monitor_login_roma",
            "value": ""
          }
    `
	helpMap["setTimeout"] = `
        - Descrizione: Imposta il tempo massimo (timeout implicito) in millisecondi che il WebDriver attenderà quando cerca un elemento (findElement, findElements) prima di restituire un errore "elemento non trovato".
        - Parametri:
          - target: Il tempo di attesa in millisecondi.
          - value: Non utilizzato.
        - Esempio:
          {
            "id": "...",
            "command": "setTimeout",
            "target": "30000", // Imposta timeout a 30 secondi
            "value": ""
          }
    `
	helpMap["activeElement"] = `
        - Descrizione: Identifica l'elemento attualmente attivo (quello con il focus) nella pagina e ne stampa i dettagli (tag, testo, stato) sulla console (output standard). Utile per debugging per capire dove si trova il focus della tastiera/interazione.
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "activeElement",
            "target": "",
            "value": ""
          }
    `
	helpMap["pageSource"] = `
        - Descrizione: Recupera l'intero sorgente HTML della pagina attualmente visualizzata e lo stampa sulla console (output standard). Utile per debugging avanzato della struttura DOM.
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "pageSource",
            "target": "",
            "value": ""
          }
    `
	helpMap["noop"] = `
        - Descrizione: Comando "No Operation". Non esegue alcuna azione. Può essere utile come segnaposto o per inserire commenti nel file .side (usando il campo comment standard di Selenium IDE, sebbene webautoma non lo legga attivamente, può essere utile per chi legge il file).
        - Parametri: Non ne usa.
        - Esempio:
          {
            "id": "...",
            "command": "noop",
            "target": "",
            "value": "",
            "comment": "Qui inizia la sezione di checkout"
          }
    `

	// Aggiungere qui altri comandi se necessario...

	return &Provider{
		commandHelp: helpMap,
	}
}

// GetCommand cerca e restituisce la stringa di aiuto per il comando specificato.
// Restituisce un errore se il comando non è trovato nella mappa di aiuto.
func (h *Provider) GetCommand(commandName string) (string, error) {
	helpText, found := h.commandHelp[commandName]
	if !found {
		return "", fmt.Errorf("nessun aiuto trovato per il comando: '%s'", commandName)
	}
	// Rimuovi spazi bianchi iniziali/finali dal blocco di testo per pulizia
	return strings.TrimSpace(helpText), nil
}

// GetSupportedCommands restituisce un elenco ordinato dei nomi dei comandi
// per i quali è disponibile l'aiuto.
func (h *Provider) GetSupportedCommands() []string {
	keys := make([]string, 0, len(h.commandHelp))
	for k := range h.commandHelp {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (h *Provider) GetFormatted(commandName string) (string, error) {
	helpText, err := h.GetCommand(commandName)
	if err != nil {
		return "", err
	}
	lines := strings.Split(helpText, "\n")
	minIndent := 999
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " ")
		if len(trimmed) == 0 {
			continue
		}
		indent := len(line) - len(trimmed)
		if indent < minIndent {
			minIndent = indent
		}
	}
	var builder strings.Builder
	for i, line := range lines {
		if len(line) >= minIndent {
			builder.WriteString(line[minIndent:])
		} else {
			builder.WriteString(line) // Mantieni linee vuote o con indentazione minore
		}
		if i < len(lines)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String(), nil
}
