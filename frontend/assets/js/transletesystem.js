 AOS.init({ duration: 800, once: true });

        // --- TRANSLATION SYSTEM ---
        const i18n = {
            ar: {
                nav_impact: "أثرنا الاجتماعي",
                nav_banks: "مراكز التبرع",
                nav_how: "كيف يعمل؟",
                nav_partners: "شركاؤنا",
                btn_login: "دخول",
                btn_donate: "تبرع الآن",
                badge_text: "مبادرة وطنية تضامنية 100%",
                hero_title_1: "تبرع بالدم..",
                hero_title_2: "أنقذ حياة مغربية.",
                hero_desc: "منصة LifeLine تربط المتبرعين بالمرضى والمستشفيات في الوقت الحقيقي. انضم إلى الآلاف من المغاربة الذين يساهمون في إنقاذ الأرواح يومياً.",
                tab_find_donor: "أبحث عن متبرع",
                tab_donate: "أريد التبرع",
                select_city: "اختر المدينة...",
                btn_search: "بحث",
                quote_text: "\"التبرع بالدم واجب وطني وإنساني\"",
                quote_author: "صاحب الجلالة الملك محمد السادس نصره الله",
                sponsors_title: "شركاؤنا في النجاح",
                steps_title: "كيفاش تبرع تتخدم؟",
                step_1_title: "1. إنشاء حساب",
                step_1_desc: "سجل ودير البروفيل ديالك، وعطي فصيلة الدم ديالك باش نلقاوك بسهولة.",
                step_2_title: "2. طلب الدم",
                step_2_desc: "محتاج دم؟ حدد الفصيلة والمكان، وحنا نتكلفو بالباقي.",
                step_3_title: "3. إشعارات ذكية",
                step_3_desc: "السيستيم كيصيفط ميساج غير للناس القراب ليك واللي عندهم نفس الزمرة.",
                step_4_title: "4. إنقاذ حياة",
                step_4_desc: "المتبرع كيوصل للمركز، كيتبرع، وكيرسم البسمة على وجه عائلة مغربية.",
                banks_title: "وضعية مخزون الدم",
                banks_subtitle: "تحديث فوري من قاعدة البيانات الوطنية",
                footer_desc: "مشروع مفتوح المصدر يهدف لرقمنة قطاع التبرع بالدم في المغرب. نحن نؤمن بأن التكنولوجيا يمكن أن تنقذ الأرواح.",
                footer_links: "روابط سريعة"
            },
            en: {
                nav_impact: "Our Impact",
                nav_banks: "Blood Banks",
                nav_how: "How it Works",
                nav_partners: "Partners",
                btn_login: "Login",
                btn_donate: "Donate Now",
                badge_text: "100% Moroccan Initiative",
                hero_title_1: "Donate Blood..",
                hero_title_2: "Save a Life.",
                hero_desc: "LifeLine connects donors, patients, and hospitals in real-time. Join thousands of Moroccans saving lives every day.",
                tab_find_donor: "Find Donor",
                tab_donate: "I Want to Donate",
                select_city: "Select City...",
                btn_search: "Search",
                quote_text: "\"Blood donation is a national and humanitarian duty\"",
                quote_author: "HM King Mohammed VI",
                sponsors_title: "Our Partners",
                steps_title: "How it Works?",
                step_1_title: "1. Create Profile",
                step_1_desc: "Register, set your blood type, and verify your details so we can reach you.",
                step_2_title: "2. Request Blood",
                step_2_desc: "Need blood? Specify the type and location, we handle the rest.",
                step_3_title: "3. Smart Alerts",
                step_3_desc: "The system notifies only nearby donors with the matching blood type.",
                step_4_title: "4. Save a Life",
                step_4_desc: "The donor arrives, donates, and brings a smile to a Moroccan family.",
                banks_title: "Blood Stock Levels",
                banks_subtitle: "Real-time updates from national database",
                footer_desc: "An open-source project digitizing blood donation in Morocco. We believe technology can save lives.",
                footer_links: "Quick Links"
            },
            fr: {
                nav_impact: "Impact",
                nav_banks: "Banques de Sang",
                nav_how: "Comment ça marche",
                nav_partners: "Partenaires",
                btn_login: "Connexion",
                btn_donate: "Donner",
                badge_text: "Initiative 100% Marocaine",
                hero_title_1: "Donnez votre sang..",
                hero_title_2: "Sauvez une vie.",
                hero_desc: "LifeLine connecte les donneurs, les patients et les hôpitaux en temps réel.",
                tab_find_donor: "Trouver un donneur",
                tab_donate: "Je veux donner",
                select_city: "Ville...",
                btn_search: "Chercher",
                quote_text: "\"Le don de sang est un devoir national et humanitaire\"",
                quote_author: "SM le Roi Mohammed VI",
                sponsors_title: "Nos Partenaires",
                steps_title: "Comment ça marche ?",
                step_1_title: "1. Créer un profil",
                step_1_desc: "Inscrivez-vous, renseignez votre groupe sanguin et vos coordonnées.",
                step_2_title: "2. Demander du sang",
                step_2_desc: "Besoin de sang ? Précisez le type et le lieu, nous gérons le reste.",
                step_3_title: "3. Alertes intelligentes",
                step_3_desc: "Le système notifie uniquement les donneurs compatibles à proximité.",
                step_4_title: "4. Sauver une vie",
                step_4_desc: "Le donneur arrive, donne son sang et sauve une vie marocaine.",
                banks_title: "Niveau de stock",
                banks_subtitle: "Mise à jour en temps réel",
                footer_desc: "Un projet open source pour la numérisation du don de sang au Maroc.",
                footer_links: "Liens Rapides"
            }
        };

        function changeLanguage(lang) {
            // 1. Update Text
            const elements = document.querySelectorAll('[data-i18n]');
            elements.forEach(el => {
                const key = el.getAttribute('data-i18n');
                if (i18n[lang][key]) {
                    el.innerText = i18n[lang][key];
                }
            });

            // 2. Update Direction & Font
            const body = document.body;
            const langLabel = document.getElementById('current-lang');
            
            if (lang === 'ar') {
                body.setAttribute('dir', 'rtl');
                body.classList.remove('font-sans');
                body.classList.add('font-arabic');
                langLabel.innerText = "العربية";
            } else {
                body.setAttribute('dir', 'ltr');
                body.classList.remove('font-arabic');
                body.classList.add('font-sans');
                langLabel.innerText = lang === 'en' ? "English" : "Français";
            }
        }
