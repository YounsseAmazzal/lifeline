Hada howa **"L-Scenario L-Kamel"** (The Full User Journey).
Mlli tbghi tchr7 l-projet, khassk t-tkhayel rassek kat7ki **9issa (Story)** fiha 4 dyal l-abtal (Actors): **L-Mrid (Patient)**, **L-Motabari3 (Donor)**, **L-Bank (Centre)**, w **L-Sponsor**.

Hahowa l-Process Step-by-Step, kifach kydour "D-dem" f l-platform:

---

### 1. Mar7alat T-tasjil (The Entry Point)
**L-Batal:** Ayman (Bgha ytbr3)
*   **L-Action:** Ayman kiy-dkhol l `index.html`. Kiy-chouf l-vibe l-malaki, kiy-ti9 (Trust).
*   **L-Process:**
    1.  Kiy-cliki "Register" (`register.html`).
    2.  Kiy-3mmar: Smia, Ville (Rabat), w **Zomra (O+)**.
    3.  **Backend Magic:** L-Go Backend kiy-chad l-address dyalo w kiy-7awlha l **GPS Coordinates** (Lat/Long) bst3mal PostGIS.
    4.  **Result:** Ayman wlla "Point Vert" f l-carte, walakin "Inactive" (Na3ess).

---

### 2. Mar7alat T-Talab (The Trigger)
**L-Batal:** Sarah (Mrid aw 3a2ila dyal mrid)
*   **L-Action:** Sarah me7taja d-dem **O+** f sbitar "Ibn Sina".
*   **L-Process:**
    1.  Kat-dkhol l `dashboard.html` (User App).
    2.  Kat-cliki "Request Blood".
    3.  **Security Check:** Kat-souwer "L-wer9a d-tbib" (Prescription) w t-uploadiha.
    4.  **Backend Magic:** L-Request mamchach l-public nichan. Mcha l **Pending Queue** 3nd l-Admin/Bank.

---

### 3. Mar7alat T-Ta7a9o9 (The Filter)
**L-Batal:** Dr. Amine (Admin aw Bank Manager)
*   **L-Action:** Kiy-jri s-sbitar mn `bank.html` aw `admin.html`.
*   **L-Process:**
    1.  Kiy-weslo notification: "New Request from Sarah".
    2.  Kiy-chouf tswira d-tbib. "Okay, hadchi s7i7".
    3.  Kiy-cliki **"Approve & Broadcast"**.
    4.  **Backend Magic (The Uber Logic):**
        *   Backend kiy-goul l Database: *"Jbed liya ga3 n-nass li O+ li saknin f rayon dyal 5km 7da Ibn Sina"*.
        *   Database kat-rad 50 wa7ed (mnhom Ayman).
        *   Backend kiy-sifet **Push Notification** l had 50 wa7ed fa9t.

---

### 4. Mar7alat L-Istijaba (The Hero Moment)
**L-Batal:** Ayman (L-Motabari3)
*   **L-Action:** Ayman jals f l-9hwa, telephon sonna.
*   **L-Process:**
    1.  Notification: *"Urgent: O+ needed 2km away. Save a life?"*
    2.  Kiy-7el `dashboard.html`. Kiy-chouf l-Map fiha Point Rouge (Sbitar) 7dah.
    3.  Kiy-cliki **"I Accept"**.
    4.  L-App kat-3tih Trajet (Map Direction) l Sbitar.

---

### 5. Mar7alat T-Tabaro3 (Closing the Loop)
**L-Batal:** Ayman & Dr. Amine
*   **L-Action:** Ayman wsal l sbitar, tbarra3.
*   **L-Process:**
    1.  Dr. Amine f `bank.html` kiy-cliki 3la profile d Ayman: **"Donation Complete"**.
    2.  **Backend Updates:**
        *   Stock dyal O+ kiy-tzad (+1 bag) f `bank.html`.
        *   Sarah (L-mrid) kat-wslha notification: "Donor found. Blood is ready."
        *   Ayman kiy-tzad lih "Points" f `profile.html` (Gamification).

---

### 6. Mar7alat L-Business (The Sponsor)
**L-Batal:** Attijariwafa Bank (Sponsor)
*   **L-Action:** Directeur Marketing kiy-dkhol l `sponsor.html`.
*   **L-Process:**
    1.  Kiy-chouf Dashboard.
    2.  Kiy-l9a: "This week, 200 lives saved using LifeLine".
    3.  Kiy-chouf Logo dyal Attijari ban 50,000 marra f App.
    4.  **Result:** Kiy-jadded l-contrat dyal sponsorship (Flous kat-dkhol l l-project).

---

### Kifach t-chr7ha b "Jomla Wa7da" (The Pitch)?

> "LifeLine hiya **Uber dyal d-dem**.
> Blast ma n-chdo taxi, kan-connectiw **Talab (Mrid)** m3a **L-3ard (Motabari3)** f l-waqt l-7aqiqi bst3mal **GPS**, w hadchi kamel **m-sécurisé** b validation dyal tbib w **m-memwel** mn charikat kbar."

### Diagramme d-Chra7 (Visual Flow):

```text
[ Mrid ] --(Talab + Tswira)--> [ Admin/Bank ] --(Validation)--> [ Backend Go ]
                                                                     |
                                                               (Calcul GPS)
                                                                     |
                                                                     v
[ Sponsor ] <--(Stats)-- [ Backend ] --(Notification)---> [ Motabari3 (Ayman) ]
                                                                     |
                                                               (Accept & Go)
                                                                     |
                                                                     v
                                                               [ Sbitar ] --(Stock Update)--> [ Fin ]
```

Hakka khassk t-chr7ha. Logic mratb w Business Model waDe7.