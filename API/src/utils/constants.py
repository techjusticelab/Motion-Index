"""
Constants used throughout the application.
"""
from typing import Dict, Set, Any

# File Processing Constants
MAX_FILE_SIZE_DEFAULT = 50 * 1024 * 1024  # 50MB

# These formats are supported by the textract library
# Reference: https://github.com/deanmalmgren/textract
SUPPORTED_FORMATS: Set[str] = {
    '.csv', '.doc', '.docx', '.eml', '.epub', '.gif', '.htm', '.html',
    '.jpeg', '.jpg', '.json', '.log', '.mp3', '.msg', '.odt', '.ogg',
    '.pdf', '.png', '.pptx', '.ps', '.psv', '.rtf', '.tab', '.tff',
    '.tif', '.tiff', '.tsv', '.txt', '.wav', '.xls', '.xlsx', 'wpd'
}

# Document Types
DOCUMENT_TYPES = {
    'motion': 'Motion',
    'petition': 'Petition',
    'order': 'Order',
    'brief': 'Brief',
    'report': 'Report',
    'exhibit': 'Exhibit',
    'memorandum': 'Memorandum',
    'response': 'Response',
    'opposition': 'Opposition',
    'complaint': 'Complaint',
    'answer': 'Answer',
    'discovery_request': 'Discovery Request',
    'discovery_response': 'Discovery Response',
    'notice': 'Notice',
    'declaration': 'Declaration',
    'affidavit': 'Affidavit',
    'judgment': 'Judgment',
    'transcript': 'Transcript',
    'settlement_agreement': 'Settlement Agreement',
    'unknown': 'Unknown'
}

# File Extension Mappings
MIME_TYPES: Dict[str, str] = {
    '.wpd': 'application/x-wordperfect',
    '.wp': 'application/x-wordperfect',
    '.wp5': 'application/x-wordperfect',
    '.mot': 'application/msword',
    '.mtn': 'application/msword',
    '.pet': 'text/plain',
    '.sup': 'text/plain',
    '.wrt': 'text/plain',
    '.reh': 'text/plain'
}

# Document Category Mappings
EXTENSION_CATEGORIES: Dict[str, str] = {
    '.pdf': 'PDF Document',
    '.docx': 'Word Document',
    '.doc': 'Word Document',
    '.wpd': 'WordPerfect Document',
    '.wp': 'WordPerfect Document',
    '.wp5': 'WordPerfect Document',
    '.txt': 'Text Document',
    '.mot': 'Motion',
    '.mtn': 'Motion',
    '.pet': 'Petition',
    '.sup': 'Supplement',
    '.ord': 'Order',
    '.rep': 'Report',
    '.ppt': 'Presentation',
    '.pptx': 'Presentation'
}

# Elasticsearch Constants
ES_DEFAULT_HOST = "localhost"
ES_DEFAULT_PORT = 9200
ES_DEFAULT_INDEX = "documents"
ES_BULK_CHUNK_SIZE = 500

# Elasticsearch Document Mapping
ES_DOCUMENT_MAPPING = {
    "mappings": {
        "properties": {
            "file_path": {"type": "keyword"},
            "file_name": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
            "category": {"type": "keyword"},
            "chunk_id": {"type": "integer"},
            "text": {"type": "text", "analyzer": "english"},
            "doc_type": {"type": "keyword"},
            "s3_uri": {"type": "keyword"},
            "metadata": {
                "properties": {
                    "document_name": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
                    "subject": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
                    "status": {"type": "keyword"},
                    "case_name": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
                    "case_number": {"type": "keyword"},
                    "author": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
                    "judge": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
                    "court": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
                    "timestamp": {"type": "date"}
                }
            },
            "embedding": {"type": "dense_vector", "dims": 384},
            "hash": {"type": "keyword"},
            "created_at": {"type": "date"}
        }
    }
}

LEGAL_TAGS_LIST = [
        "Abandonment of Criminal Purpose", "Abatement", "Accessory After the Fact", "Accomplice Liability",
        "Acquittal", "Actus Reus", "Administrative Law Violations", "Administrative Search",
        "Admissibility of Evidence", "Affirmative Defenses", "Age of Consent", "Aggravated Assault",
        "Aggravated Battery", "Aggravating Factors", "Aiding and Abetting", "Alcohol-Related Offenses",
        "Alibi Defense", "Alternate Jurors", "Amendment of Charges", "Animal Cruelty",
        "Anonymous Tips", "Anti-Terrorism Laws", "Appeals", "Appellate Jurisdiction",
        "Arrest", "Arrest Record Expungement", "Arrest Warrant", "Arson",
        "Assault and Battery", "Assumption of Risk", "Attempt Crimes", "Attorney-Client Privilege",
        "Automobile Exception", "Autopsy Evidence", "Bail and Pretrial Release", "Ballistics Evidence",
        "Banking Crimes", "Bench Trial", "Bias-Motivated Crimes", "Bifurcated Trial",
        "Bill of Particulars", "Blackmail", "Blood Alcohol Content", "Body-Worn Camera Evidence",
        "Bond Forfeiture", "Border Searches", "Brady Violations", "Brandishing a Weapon",
        "Bribery", "Burden of Proof", "Burglary", "Business Records Exception",
        "Capital Punishment", "Carjacking", "Chain of Custody", "Change of Venue",
        "Character Evidence", "Charging Documents", "Check Fraud", "Child Abuse",
        "Child Endangerment", "Child Molestation", "Child Pornography", "Circumstantial Evidence",
        "Civil Commitment", "Civil Forfeiture", "Clergy-Penitent Privilege", "Co-Conspirator Statements",
        "Coerced Confession", "Collateral Attack", "Commutation of Sentence", "Community Property Crimes",
        "Community Service", "Competency to Stand Trial", "Complaint", "Computer Crimes",
        "Concealed Weapons", "Concurrent Jurisdiction", "Concurrent Sentences", "Conditional Plea",
        "Confessions", "Confidential Informants", "Confrontation Clause", "Consent Searches",
        "Conspiracy", "Constitutional Law", "Constructive Possession", "Contempt of Court",
        "Continuance", "Contraband", "Controlled Substances", "Coroner's Report",
        "Corporate Criminal Liability", "Corpus Delicti", "Correctional Facilities", "Corroboration Requirement",
        "Counterfeiting", "Court-Appointed Counsel", "Credit Card Fraud", "Credit for Time Served",
        "Criminal Enterprise", "Criminal Forfeiture", "Criminal History", "Criminal Negligence",
        "Criminal Procedure", "Criminal Registration", "Criminal Solicitation", "Cross-Examination",
        "Cruel and Unusual Punishment", "Culpability", "Cultural Defense", "Cumulative Punishment",
        "Custodial Interrogation", "Custody", "Cyberbullying", "Cyberstalking",
        "Deadly Weapon Enhancement", "Death Penalty", "Defamation", "Default Judgment",
        "Defective Warrant", "Defenses", "Deferred Adjudication", "Deferred Sentence",
        "Deliberate Indifference", "Deliberations", "Delinquency Proceedings", "Denial of Counsel",
        "Deportation Consequences", "Depositions", "Destruction of Evidence", "Detention",
        "Determinate Sentencing", "Diminished Capacity", "Direct Examination", "Direct Evidence",
        "Direct File (Juvenile)", "Discovery", "Dismissal With Prejudice", "Dismissal Without Prejudice",
        "Disorderly Conduct", "Diversion Programs", "DNA Evidence", "Doctrine of Chances",
        "Doctrine of Transferred Intent", "Domestic Violence", "Double Jeopardy", "Drug Cultivation",
        "Drug Distribution", "Drug Manufacturing", "Drug Offenses", "Drug Paraphernalia",
        "Drug Testing", "Due Process", "Duress Defense", "DUI/DWI",
        "Elder Abuse", "Electronic Surveillance", "Elements of Crimes", "Embezzlement",
        "Emergency Exception", "Emotional Distress", "En Banc Review", "Enhancement Allegations",
        "Entrapment", "Environmental Crimes", "Equal Protection", "Escape",
        "Evading Arrest", "Evidence", "Evidence Tampering", "Ex Parte Communications",
        "Ex Post Facto Laws", "Excessive Bail", "Excessive Force", "Exclusionary Rule",
        "Exculpatory Evidence", "Excusable Homicide", "Expert Testimony", "Extortion",
        "Extradition", "Eyewitness Identification", "Factual Impossibility", "Failure to Appear",
        "False Arrest", "False Imprisonment", "False Statements", "Family Violence",
        "Federal Preemption", "Felonies", "Felony Murder", "Fiduciary Duty Violations",
        "Field Sobriety Tests", "Fifth Amendment", "Financial Crimes", "Fingerprint Evidence",
        "Firearm Enhancement", "Firearm Possession", "First Amendment Issues", "Flight as Evidence",
        "Forfeiture", "Forgery", "Fourth Amendment", "Fraud",
        "Fresh Pursuit", "Fruit of the Poisonous Tree", "Fugitive Status", "Functus Officio",
        "Gang Enhancements", "Good Faith Exception", "Grand Jury", "Great Bodily Injury",
        "Habeas Corpus", "Habitual Offender Status", "Harmless Error", "Hate Crimes",
        "Hearsay", "Hearsay Exceptions", "Hit and Run", "Home Detention",
        "Homicide", "Honeypot Operations", "Hot Pursuit", "Human Trafficking",
        "Identity Theft", "Illegal Gambling", "Illegal Search", "Immunity",
        "Impeachment Evidence", "Implied Consent", "Impoundment", "Improperly Joined Charges",
        "Incarceration", "Inchoate Crimes", "Indecent Exposure", "Indeterminate Sentencing",
        "Indictment", "Indigent Defense", "Ineffective Assistance of Counsel", "Information (Charging Document)",
        "Initial Appearance", "Injunctions", "Insanity Defense", "Instruction Conference",
        "Insufficient Evidence", "Insurance Fraud", "Intent", "Intent to Distribute",
        "Interlocutory Appeal", "Internal Affairs Investigation", "International Criminal Law", "Internet Crimes",
        "Interrogation", "Inventory Search", "Involuntary Intoxication", "Involuntary Manslaughter",
        "Jailhouse Informants", "Joint Trial", "Judicial Bias", "Judicial Notice",
        "Jurisdiction", "Jury Instructions", "Jury Nullification", "Jury Selection",
        "Justifiable Homicide", "Juvenile Adjudication", "Juvenile Law", "Kidnapping",
        "Knock and Announce", "Knowledge Element", "Laboratory Reports", "Larceny",
        "Law Enforcement Testimony", "Legal Malpractice", "Lesser Included Offenses", "Lex Talionis",
        "Liability Without Fault", "Libel", "Lineup Procedures", "Loitering",
        "Lone Actor Terrorism", "Loyalty Crimes", "Mail Fraud", "Malice",
        "Mandatory Minimum Sentences", "Mandatory Reporting", "Manslaughter", "Material Witness",
        "Mayhem", "Mens Rea", "Mental Health Diversion", "Mental State",
        "Military Justice", "Ministerial Acts", "Miranda Rights", "Misdemeanors",
        "Misprision of Felony", "Missing Evidence", "Mitigating Factors", "Mitigating Circumstances",
        "Modified Categorical Approach", "Money Laundering", "Mootness", "Motion for Acquittal",
        "Motion for New Trial", "Motion to Dismiss", "Motion to Quash", "Motion to Suppress",
        "Motions", "Multiple Punishment", "Murder", "National Security Exception",
        "Necessity Defense", "Negligent Homicide", "No Contest Plea", "Nolo Contendere",
        "Non-Testimonial Evidence", "Nolle Prosequi", "Nondisclosure Agreement Violations", "Objectivity",
        "Obscenity", "Obstruction of Justice", "Officer Safety Exception", "Official Immunity",
        "Omnibus Hearing", "Open Fields Doctrine", "Operating a Vehicle Under the Influence", "Opinion Testimony",
        "Opportunity to Commit", "Oral Arguments", "Ordinance Violations", "Overbreadth",
        "Overcharging", "Overt Act", "Pardon", "Parole",
        "Parole Evidence Rule", "Participants in Crime", "Particularity Requirement", "Party Admissions",
        "Pattern of Criminal Activity", "Peace Officer Status", "Penal Code Violations", "Penitentiary Time",
        "Per Se Violations", "Peremptory Challenges", "Perjury", "Permissive Inference",
        "Personal Jurisdiction", "Petition for Writ", "Photographic Evidence", "Physical Evidence",
        "Plain Error", "Plain Feel Doctrine", "Plain View Doctrine", "Plea Bargaining",
        "Posse Comitatus", "Post-Conviction Relief", "Post-Release Community Supervision", "Precedential Value",
        "Preemptive Self-Defense", "Preliminary Hearings", "Premeditation", "Preponderance of Evidence",
        "Prescription Drug Fraud", "Presentence Investigation", "Presentment", "Presumption of Innocence",
        "Presumptions", "Pretextual Stops", "Pretrial Conference", "Pretrial Diversion",
        "Pretrial Identification", "Pretrial Motions", "Pretrial Release Conditions", "Prima Facie Case",
        "Prior Bad Acts", "Prior Conviction Enhancement", "Prison Conditions", "Privacy Rights",
        "Private Search Doctrine", "Privilege Against Self-Incrimination", "Probable Cause", "Probation",
        "Probation Violation", "Procedural Due Process", "Professional Responsibility", "Proof Beyond Reasonable Doubt",
        "Property Crimes", "Proportionality Review", "Prosecution", "Prosecutorial Discretion",
        "Prosecutorial Immunity", "Prosecutorial Misconduct", "Protected Classes", "Protective Orders",
        "Proximate Cause", "Psychiatric Examination", "Public Authority Defense", "Public Defender",
        "Public Duty Doctrine", "Public Necessity", "Public Safety Exception", "Public Trial",
        "Public Trust Crimes", "Punishment", "Pyramid Schemes", "Qualified Immunity",
        "Racketeering", "Rape and Sexual Assault", "Rape Shield Laws", "Ratification",
        "Reasonable Doubt", "Reasonable Force", "Reasonable Person Standard", "Reasonable Suspicion",
        "Receiving Stolen Property", "Reckless Driving", "Recklessness", "Recording of Interrogation",
        "Recusal", "Reduced Charges", "Regulatory Offenses", "Rehabilitation",
        "Relevant Evidence", "Remedial Measures", "Remittitur", "Reopening Case",
        "Repeated Offenses", "Required Records Doctrine", "Res Gestae", "Res Judicata",
        "Resisting Arrest", "Restitution", "Restraining Orders", "Resultant Crimes",
        "Retroactive Application", "Retroactivity", "Return of Property", "Revenue Crimes",
        "Reverse Sting Operations", "Revocation Hearing", "Right of Confrontation", "Right to a Speedy Trial",
        "Right to Counsel", "Right to Present Evidence", "Right to Remain Silent", "Risk Assessment",
        "Robbery", "Rules of Criminal Procedure", "Rules of Evidence", "Safety Valve Provisions",
        "Sanctions", "Scientific Evidence", "Scope of Consent", "Scope of Search",
        "Search and Seizure", "Search Incident to Arrest", "Search Warrant", "Secondary Evidence",
        "Secret Witness", "Security Interests", "Sedition", "Self-Defense",
        "Self-Incrimination", "Sentencing", "Sentencing Enhancement", "Sentencing Guidelines",
        "Sentencing Hearing", "Sequential Lineup", "Sex Offender Registration", "Sexual Battery",
        "Sexual Exploitation", "Sexual Predator Laws", "Sexually Violent Predator", "Shackling Defendants",
        "Shadow Jury", "Shoplifting", "Show-Up Identification", "Signature Crime",
        "Similar Transactions", "Simulated Controlled Substances", "Single Transaction Rule", "Sixth Amendment Rights",
        "Slander", "Sneak and Peek Warrant", "Sobriety Checkpoint", "Solicitation",
        "Special Allegations", "Special Circumstances", "Special Findings", "Special Relationship",
        "Specific Intent", "Speedy Trial", "Split Sentence", "Spoliation of Evidence",
        "Stalking", "Standing", "State Action Doctrine", "Statement Against Interest",
        "Statute of Limitations", "Statutory Interpretation", "Statutory Rape", "Statutory Remedies",
        "Stop and Frisk", "Straw Purchase", "Strict Liability", "Strict Scrutiny",
        "Strike Offenses", "Subpoenas", "Substantial Compliance", "Substantial Evidence",
        "Substantive Due Process", "Substitution of Counsel", "Suggestive Identification", "Suicide by Cop",
        "Summary Judgment", "Summons", "Supervised Release", "Suppression Hearing",
        "Suppression of Evidence", "Sua Sponte", "Suspended Sentence", "Sustaining Objections",
        "Sweet Plea", "Tangible Evidence", "Taser Use", "Tax Evasion",
        "Technical Violations", "Temporary Detention", "Territorial Jurisdiction", "Terrorism",
        "Testimonial Evidence", "Testimony", "Theft", "Third Party Culpability",
        "Three Strikes Law", "Time Credit", "Time Limits", "Title IX Violations",
        "Torts and Crimes", "Traffic Offenses", "Transferred Intent", "Transitional Release",
        "Transportation of Contraband", "Trap and Trace Devices", "Trespass to Investigate", "Trial by Declaration",
        "Trial Court Discretion", "Trial Procedures", "Trial Strategy", "Trucking Violations",
        "True Bill", "Truncated Investigation", "Unconstitutional Conditions", "Undercover Operations",
        "Undue Influence", "Unlawful Assembly", "Unlawful Detention", "Unlawful Entry",
        "Unlawful Search", "Unlicensed Practice", "Use Immunity", "Uttering",
        "Vagueness Challenge", "Vandalism", "Vehicle Code Violations", "Vehicle Exception",
        "Vehicle Searches", "Vehicle Stop", "Venue", "Verdict Forms",
        "Verdicts", "Victim Advocate", "Victim Impact Evidence", "Victim Impact Statement",
        "Victim Rights", "Victim-Witness Protection", "Video Evidence", "Violation of Court Order",
        "Violent Felony", "Voir Dire", "Voluntary Act", "Voluntary Intoxication",
        "Voluntary Manslaughter", "Waiver of Appeal", "Waiver of Counsel", "Waiver of Miranda Rights",
        "Waiver of Preliminary Hearing", "Waiver of Rights", "Warrantless Arrest", "Warrantless Search",
        "Warrants", "Weapons Enhancement", "Weapons Possession", "White Collar Crime",
        "Willful Blindness", "Wire Fraud", "Wiretapping", "Withdrawal from Conspiracy",
        "Witness Credibility", "Witness Impeachment", "Witness Intimidation", "Witness Testimony",
        "Wolf Pack Prosecution", "Work Release", "Workplace Violence", "Wrongful Conviction",
        "Wrongful Death", "Youthful Offender", "Zero Tolerance Policies", "Zoning Violations",
    ]

# Processing Constants
DEFAULT_MAX_WORKERS = 4
DEFAULT_BATCH_SIZE = 100

# LLM Constants
OPENAI_MODEL = "gpt-3.5-turbo"
