Attaching to app-1
app-1  | 2025-05-01 00:21:55,053 - elasticsearch_handler - INFO - Using Elasticsearch at my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443
app-1  | 2025-05-01 00:21:55,054 - elasticsearch_handler - INFO - Using API key authentication
app-1  | 2025-05-01 00:21:55,937 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/ [status:200 duration:0.864s]
app-1  | 2025-05-01 00:21:55,937 - elasticsearch_handler - INFO - Successfully connected to Elasticsearch
app-1  | 2025-05-01 00:21:56,063 - s3_handler - INFO - Connected to S3 with bucket: cpda-documents
app-1  | 2025-05-01 00:21:56,346 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda [status:200 duration:0.282s]
app-1  | 2025-05-01 00:21:56,346 - elasticsearch_handler - INFO - Index 'cpda' already exists
app-1  | 2025-05-01 00:21:56,346 - motion_index - INFO - Starting document processing for directory: /data
app-1  | 2025-05-01 00:21:56,346 - document_processor - INFO - Processing directory: /data
app-1  | 2025-05-01 00:21:56,346 - document_processor - WARNING - File extension  not supported by textract: /data/1538motwpd
app-1  | 2025-05-01 00:21:56,346 - document_processor - WARNING - File extension .mot not supported by textract: /data/1538.mot
app-1  | 2025-05-01 00:21:56,347 - document_processor - INFO - Found 17 total files
app-1  | 2025-05-01 00:21:56,347 - document_processor - INFO - Skipped 2 files (too large: 0, unsupported format: 2)
app-1  | 2025-05-01 00:21:56,347 - document_processor - INFO - Will process 15 files
Processing files:   0%|          | 0/15 [00:00<?, ?it/s]2025-05-01 00:21:56,494 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/f7cdc489d8b31315b1987026cc028edd [status:404 duration:0.136s]
app-1  | 2025-05-01 00:21:56,694 - document_processor - INFO - Classifying document: 1538FORM.pdf
app-1  | 2025-05-01 00:21:56,695 - document_classifier - INFO - Attempting LLM classification and extraction...
app-1  | 2025-05-01 00:21:56,755 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/c3ce0862adb0db4d235dfb941fba6bfa [status:404 duration:0.395s]
app-1  | 2025-05-01 00:21:56,756 - file_processor - INFO - Converting WPD file to PDF: /data/995_pd.wpd
app-1  | 2025-05-01 00:21:56,756 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/958bc61f71fbfcdc6c51c3ab507f9854 [status:404 duration:0.396s]
app-1  | 2025-05-01 00:21:56,757 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/fb99c7ecf4fc3d299ecd3c335d74c2fc [status:404 duration:0.397s]
app-1  | 2025-05-01 00:21:56,766 - file_processor - INFO - Using temp directory: /tmp/tmpcq2ibqwf
app-1  | 2025-05-01 00:21:56,766 - file_processor - INFO - Base filename: 995_pd.wpd
app-1  | 2025-05-01 00:21:56,766 - file_processor - INFO - Expected PDF path: /tmp/tmpcq2ibqwf/995_pd.pdf
app-1  | 2025-05-01 00:21:57,035 - file_processor - INFO - Found LibreOffice at: libreoffice
app-1  | 2025-05-01 00:21:57,035 - file_processor - INFO - Running conversion command: libreoffice --headless --convert-to pdf --outdir /tmp/tmpcq2ibqwf /data/995_pd.wpd
app-1  | 2025-05-01 00:21:57,170 - document_processor - INFO - Classifying document: admit_~2.pdf
app-1  | 2025-05-01 00:21:57,170 - document_classifier - INFO - Attempting LLM classification and extraction...
app-1  | 2025-05-01 00:21:58,211 - file_processor - INFO - Command stdout:
app-1  | 2025-05-01 00:21:58,211 - file_processor - INFO - Command stderr: Warning: failed to launch javaldx - java may not function correctly
app-1  | Error: source file could not be loaded
app-1  |
app-1  | 2025-05-01 00:21:58,211 - file_processor - INFO - Checking output directory contents: /tmp/tmpcq2ibqwf
app-1  | 2025-05-01 00:21:58,211 - file_processor - INFO - Files in output directory: []
app-1  | 2025-05-01 00:21:58,211 - file_processor - ERROR - No PDF files found in output directory
app-1  | 2025-05-01 00:21:58,211 - file_processor - ERROR - Failed to convert WPD file to PDF: /data/995_pd.wpd
app-1  | 2025-05-01 00:21:58,211 - document_processor - WARNING - No text extracted from /data/995_pd.wpd
Processing files:   7%|▋         | 1/15 [00:01<00:25,  1.86s/it]2025-05-01 00:21:58,393 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/6dbe572104cafcf9ede6883ce51d6de2 [status:404 duration:0.181s]
app-1  | 2025-05-01 00:21:58,393 - file_processor - INFO - Converting WPD file to PDF: /data/1538mot.wpd
app-1  | 2025-05-01 00:21:58,394 - file_processor - INFO - Using temp directory: /tmp/tmpkwxutqf9
app-1  | 2025-05-01 00:21:58,394 - file_processor - INFO - Base filename: 1538mot.wpd
app-1  | 2025-05-01 00:21:58,394 - file_processor - INFO - Expected PDF path: /tmp/tmpkwxutqf9/1538mot.pdf
app-1  | 2025-05-01 00:21:58,493 - file_processor - INFO - Found LibreOffice at: libreoffice
app-1  | 2025-05-01 00:21:58,493 - file_processor - INFO - Running conversion command: libreoffice --headless --convert-to pdf --outdir /tmp/tmpkwxutqf9 /data/1538mot.wpd
app-1  | 2025-05-01 00:21:59,371 - file_processor - INFO - Command stdout:
app-1  | 2025-05-01 00:21:59,371 - file_processor - INFO - Command stderr: Warning: failed to launch javaldx - java may not function correctly
app-1  | Error: source file could not be loaded
app-1  |
app-1  | 2025-05-01 00:21:59,371 - file_processor - INFO - Checking output directory contents: /tmp/tmpkwxutqf9
app-1  | 2025-05-01 00:21:59,371 - file_processor - INFO - Files in output directory: []
app-1  | 2025-05-01 00:21:59,371 - file_processor - ERROR - No PDF files found in output directory
app-1  | 2025-05-01 00:21:59,371 - file_processor - ERROR - Failed to convert WPD file to PDF: /data/1538mot.wpd
app-1  | 2025-05-01 00:21:59,371 - document_processor - WARNING - No text extracted from /data/1538mot.wpd
Processing files:  13%|█▎        | 2/15 [00:03<00:18,  1.45s/it]2025-05-01 00:21:59,451 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/249ca9b9de2b4a8b0d011c0785334213 [status:404 duration:0.078s]
app-1  | 2025-05-01 00:21:59,739 - document_classifier - INFO - LLM classification/extraction successful: Type 'Order', Tags: ['Search and Seizure', 'Evidence', 'Suppression of Evidence', 'Rules of Criminal Procedure'], Timestamp: 'None'
app-1  | 2025-05-01 00:21:59,739 - document_processor - INFO - Document classified as: Order, Timestamp: None, Legal Tags: ['Search and Seizure', 'Evidence', 'Suppression of Evidence', 'Rules of Criminal Procedure']
app-1  | 2025-05-01 00:22:00,094 - document_classifier - INFO - LLM classification/extraction successful: Type 'Notice', Tags: ['Search and Seizure', 'Warrantless Search', 'Motion to Suppress', 'Probable Cause', 'Penal Code Violations'], Timestamp: '1999-00-00'
app-1  | 2025-05-01 00:22:00,095 - document_processor - INFO - Document classified as: Notice, Timestamp: None, Legal Tags: ['Search and Seizure', 'Warrantless Search', 'Motion to Suppress', 'Probable Cause', 'Penal Code Violations']
app-1  | 2025-05-01 00:22:00,503 - s3_handler - INFO - Uploaded /data/admit_~2.pdf to s3://cpda-documents/order/2025/05/01/admit_~2_a02f8bc1.pdf
app-1  | 2025-05-01 00:22:00,503 - document_processor - INFO - Uploaded to S3: s3://cpda-documents/order/2025/05/01/admit_~2_a02f8bc1.pdf
Processing files:  20%|██        | 3/15 [00:04<00:15,  1.30s/it]2025-05-01 00:22:00,512 - s3_handler - INFO - Uploaded /data/1538FORM.pdf to s3://cpda-documents/notice/2025/05/01/1538FORM_e3f26a2b.pdf
app-1  | 2025-05-01 00:22:00,513 - document_processor - INFO - Uploaded to S3: s3://cpda-documents/notice/2025/05/01/1538FORM_e3f26a2b.pdf
app-1  | 2025-05-01 00:22:00,588 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/715a1e7f9c2e3dcba9000904cc9c77ae [status:404 duration:0.079s]
app-1  | 2025-05-01 00:22:00,588 - file_processor - INFO - Converting WPD file to PDF: /data/1800~1.wpd
app-1  | 2025-05-01 00:22:00,588 - file_processor - INFO - Using temp directory: /tmp/tmp2upbs1bm
app-1  | 2025-05-01 00:22:00,588 - file_processor - INFO - Base filename: 1800~1.wpd
app-1  | 2025-05-01 00:22:00,588 - file_processor - INFO - Expected PDF path: /tmp/tmp2upbs1bm/1800~1.pdf
app-1  | 2025-05-01 00:22:00,593 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/9721be52756ccb38e7b8391a9b89009a [status:404 duration:0.079s]
app-1  | 2025-05-01 00:22:00,692 - file_processor - INFO - Found LibreOffice at: libreoffice
app-1  | 2025-05-01 00:22:00,692 - file_processor - INFO - Running conversion command: libreoffice --headless --convert-to pdf --outdir /tmp/tmp2upbs1bm /data/1800~1.wpd
app-1  | 2025-05-01 00:22:00,909 - document_processor - INFO - Classifying document: 98_rep~1.pdf
app-1  | 2025-05-01 00:22:00,910 - document_classifier - INFO - Attempting LLM classification and extraction...
app-1  | 2025-05-01 00:22:01,168 - document_processor - INFO - Classifying document: 2_fina~1.pdf
app-1  | 2025-05-01 00:22:01,174 - document_classifier - INFO - Attempting LLM classification and extraction...
app-1  | 2025-05-01 00:22:01,451 - document_processor - INFO - Classifying document: 00Speedy.pdf
app-1  | 2025-05-01 00:22:01,451 - document_classifier - INFO - Attempting LLM classification and extraction...
app-1  | 2025-05-01 00:22:01,622 - file_processor - INFO - Command stdout:
app-1  | 2025-05-01 00:22:01,622 - file_processor - INFO - Command stderr: Warning: failed to launch javaldx - java may not function correctly
app-1  | Error: source file could not be loaded
app-1  |
app-1  | 2025-05-01 00:22:01,622 - file_processor - INFO - Checking output directory contents: /tmp/tmp2upbs1bm
app-1  | 2025-05-01 00:22:01,622 - file_processor - INFO - Files in output directory: []
app-1  | 2025-05-01 00:22:01,622 - file_processor - ERROR - No PDF files found in output directory
app-1  | 2025-05-01 00:22:01,622 - file_processor - ERROR - Failed to convert WPD file to PDF: /data/1800~1.wpd
app-1  | 2025-05-01 00:22:01,622 - document_processor - WARNING - No text extracted from /data/1800~1.wpd
Processing files:  33%|███▎      | 5/15 [00:05<00:08,  1.13it/s]2025-05-01 00:22:01,702 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/8fe400e53e58c434407b806ec8a00f4b [status:404 duration:0.078s]
app-1  | 2025-05-01 00:22:01,986 - document_processor - INFO - Classifying document: admit_~1.pdf
app-1  | 2025-05-01 00:22:01,986 - document_classifier - INFO - Attempting LLM classification and extraction...
app-1  | 2025-05-01 00:22:03,923 - document_classifier - INFO - LLM classification/extraction successful: Type 'Report', Tags: ['Wiretapping', 'Search and Seizure', 'Electronic Surveillance', 'Probable Cause'], Timestamp: '1998-01-07'
app-1  | 2025-05-01 00:22:03,923 - document_processor - INFO - Document classified as: Report, Timestamp: 1998-01-07 00:00:00, Legal Tags: ['Wiretapping', 'Search and Seizure', 'Electronic Surveillance', 'Probable Cause']
app-1  | 2025-05-01 00:22:04,644 - document_classifier - INFO - LLM classification/extraction successful: Type 'Memorandum', Tags: ['Speedy Trial', 'Penal Code Violations', 'Due Process', 'Sixth Amendment Rights', 'Constitutional Law'], Timestamp: '2000-05-15'
app-1  | 2025-05-01 00:22:04,644 - document_processor - INFO - Document classified as: Memorandum, Timestamp: 2000-05-15 00:00:00, Legal Tags: ['Speedy Trial', 'Penal Code Violations', 'Due Process', 'Sixth Amendment Rights', 'Constitutional Law']
app-1  | 2025-05-01 00:22:04,749 - s3_handler - INFO - Uploaded /data/98_rep~1.pdf to s3://cpda-documents/report/2025/05/01/98_rep~1_5706985e.pdf
app-1  | 2025-05-01 00:22:04,749 - document_processor - INFO - Uploaded to S3: s3://cpda-documents/report/2025/05/01/98_rep~1_5706985e.pdf
Processing files:  40%|████      | 6/15 [00:08<00:13,  1.53s/it]2025-05-01 00:22:04,951 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/b72fcefc94be231450451db504d634b0 [status:404 duration:0.196s]
app-1  | 2025-05-01 00:22:05,361 - document_processor - INFO - Classifying document: 98_12.pdf
app-1  | 2025-05-01 00:22:05,361 - document_classifier - INFO - Attempting LLM classification and extraction...
app-1  | 2025-05-01 00:22:05,564 - document_classifier - INFO - LLM classification/extraction successful: Type 'Affidavit', Tags: ['Wiretapping', 'Search and Seizure', 'Suppression of Evidence', 'Procedural Due Process'], Timestamp: '1999-09-01'
app-1  | 2025-05-01 00:22:05,564 - document_processor - INFO - Document classified as: Affidavit, Timestamp: 1999-09-01 00:00:00, Legal Tags: ['Wiretapping', 'Search and Seizure', 'Suppression of Evidence', 'Procedural Due Process']
app-1  | 2025-05-01 00:22:06,178 - document_classifier - INFO - LLM classification/extraction successful: Type 'Order', Tags: ['Search and Seizure', 'Wiretapping', 'Admissibility of Evidence', 'Affirmative Defenses', 'Fourth Amendment'], Timestamp: '2000-04-03'
app-1  | 2025-05-01 00:22:06,178 - document_processor - INFO - Document classified as: Order, Timestamp: 2000-04-03 00:00:00, Legal Tags: ['Search and Seizure', 'Wiretapping', 'Admissibility of Evidence', 'Affirmative Defenses', 'Fourth Amendment']
app-1  | 2025-05-01 00:22:06,754 - s3_handler - INFO - Uploaded /data/2_fina~1.pdf to s3://cpda-documents/affidavit/2025/05/01/2_fina~1_b8f13ebf.pdf
app-1  | 2025-05-01 00:22:06,754 - document_processor - INFO - Uploaded to S3: s3://cpda-documents/affidavit/2025/05/01/2_fina~1_b8f13ebf.pdf
Processing files:  47%|████▋     | 7/15 [00:10<00:13,  1.67s/it]2025-05-01 00:22:06,830 - s3_handler - INFO - Uploaded /data/00Speedy.pdf to s3://cpda-documents/memorandum/2025/05/01/00Speedy_5e1ff812.pdf
app-1  | 2025-05-01 00:22:06,830 - document_processor - INFO - Uploaded to S3: s3://cpda-documents/memorandum/2025/05/01/00Speedy_5e1ff812.pdf
app-1  | 2025-05-01 00:22:06,839 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/8b59bcf521ac2c5f185b956f208f0295 [status:404 duration:0.078s]
app-1  | 2025-05-01 00:22:06,907 - s3_handler - INFO - Uploaded /data/admit_~1.pdf to s3://cpda-documents/order/2025/05/01/admit_~1_60c03bdb.pdf
app-1  | 2025-05-01 00:22:06,908 - document_processor - INFO - Uploaded to S3: s3://cpda-documents/order/2025/05/01/admit_~1_60c03bdb.pdf
Processing files:  60%|██████    | 9/15 [00:10<00:05,  1.06it/s]2025-05-01 00:22:06,912 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/c4cbf13a8c9b66abe35ae57a78a7020b [status:200 duration:0.080s]
app-1  | 2025-05-01 00:22:06,912 - document_processor - INFO - Skipping existing document: /data/99_0215.pdf
app-1  | 2025-05-01 00:22:06,988 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/804c83c6a32fef707f21ef16792422f2 [status:404 duration:0.078s]
app-1  | 2025-05-01 00:22:06,988 - file_processor - INFO - Converting WPD file to PDF: /data/402.wpd
app-1  | 2025-05-01 00:22:06,989 - file_processor - INFO - Using temp directory: /tmp/tmp035u03a6
app-1  | 2025-05-01 00:22:06,989 - file_processor - INFO - Base filename: 402.wpd
app-1  | 2025-05-01 00:22:06,989 - file_processor - INFO - Expected PDF path: /tmp/tmp035u03a6/402.pdf
app-1  | 2025-05-01 00:22:06,994 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/f2a2ec410b8fe7657ab1a5c24e9a6618 [status:404 duration:0.081s]
app-1  | 2025-05-01 00:22:07,128 - file_processor - INFO - Found LibreOffice at: libreoffice
app-1  | 2025-05-01 00:22:07,128 - file_processor - INFO - Running conversion command: libreoffice --headless --convert-to pdf --outdir /tmp/tmp035u03a6 /data/402.wpd
app-1  | 2025-05-01 00:22:07,174 - document_processor - INFO - Classifying document: MEMORANDUM OF UNDERSTANDING BETWEEN LACPD BERKELEY (CLJC edits).pdf
app-1  | 2025-05-01 00:22:07,174 - document_classifier - INFO - Attempting LLM classification and extraction...
app-1  | 2025-05-01 00:22:07,307 - document_processor - INFO - Classifying document: 245usen.pdf
app-1  | 2025-05-01 00:22:07,307 - document_classifier - INFO - Attempting LLM classification and extraction...
app-1  | 2025-05-01 00:22:07,775 - document_classifier - INFO - LLM classification/extraction successful: Type 'Memorandum', Tags: ['Three Strikes Law', 'Sentencing', 'Prior Bad Acts', 'Mandatory Minimum Sentences', 'Judicial Bias'], Timestamp: 'None'
app-1  | 2025-05-01 00:22:07,775 - document_processor - INFO - Document classified as: Memorandum, Timestamp: None, Legal Tags: ['Three Strikes Law', 'Sentencing', 'Prior Bad Acts', 'Mandatory Minimum Sentences', 'Judicial Bias']
app-1  | 2025-05-01 00:22:08,116 - s3_handler - INFO - Uploaded /data/98_12.pdf to s3://cpda-documents/memorandum/2025/05/01/98_12_4f3460e7.pdf
app-1  | 2025-05-01 00:22:08,116 - document_processor - INFO - Uploaded to S3: s3://cpda-documents/memorandum/2025/05/01/98_12_4f3460e7.pdf
Processing files:  73%|███████▎  | 11/15 [00:11<00:03,  1.23it/s]2025-05-01 00:22:08,202 - elastic_transport.transport - INFO - HEAD https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/cpda/_doc/f42f754888b430da9c656ded91c18fb7 [status:404 duration:0.079s]
app-1  | 2025-05-01 00:22:08,350 - file_processor - INFO - Command stdout:
app-1  | 2025-05-01 00:22:08,350 - file_processor - INFO - Command stderr: Warning: failed to launch javaldx - java may not function correctly
app-1  | Error: source file could not be loaded
app-1  |
app-1  | 2025-05-01 00:22:08,350 - file_processor - INFO - Checking output directory contents: /tmp/tmp035u03a6
app-1  | 2025-05-01 00:22:08,350 - file_processor - INFO - Files in output directory: []
app-1  | 2025-05-01 00:22:08,350 - file_processor - ERROR - No PDF files found in output directory
app-1  | 2025-05-01 00:22:08,351 - file_processor - ERROR - Failed to convert WPD file to PDF: /data/402.wpd
app-1  | 2025-05-01 00:22:08,351 - document_processor - WARNING - No text extracted from /data/402.wpd
Processing files:  80%|████████  | 12/15 [00:11<00:02,  1.46it/s]2025-05-01 00:22:08,453 - document_processor - INFO - Classifying document: 707(B).pdf
app-1  | 2025-05-01 00:22:08,454 - document_classifier - INFO - Attempting LLM classification and extraction...
app-1  | 2025-05-01 00:22:09,890 - document_classifier - INFO - LLM classification/extraction successful: Type 'Memorandum', Tags: ['Strike Offenses', 'Rules of Evidence', 'Motion to Dismiss', 'Burden of Proof'], Timestamp: 'None'
app-1  | 2025-05-01 00:22:09,891 - document_processor - INFO - Document classified as: Memorandum, Timestamp: None, Legal Tags: ['Strike Offenses', 'Rules of Evidence', 'Motion to Dismiss', 'Burden of Proof']
app-1  | 2025-05-01 00:22:10,186 - document_classifier - INFO - LLM classification/extraction successful: Type 'Memorandum', Tags: ['Military Justice', 'Attorney-Client Privilege', 'Legal Malpractice', 'Public Defender'], Timestamp: 'None'
app-1  | 2025-05-01 00:22:10,186 - document_processor - INFO - Document classified as: Memorandum, Timestamp: None, Legal Tags: ['Military Justice', 'Attorney-Client Privilege', 'Legal Malpractice', 'Public Defender']
app-1  | 2025-05-01 00:22:10,189 - s3_handler - INFO - Uploaded /data/245usen.pdf to s3://cpda-documents/memorandum/2025/05/01/245usen_767b1057.pdf
app-1  | 2025-05-01 00:22:10,189 - document_processor - INFO - Uploaded to S3: s3://cpda-documents/memorandum/2025/05/01/245usen_767b1057.pdf
Processing files:  87%|████████▋ | 13/15 [00:13<00:01,  1.04it/s]2025-05-01 00:22:10,417 - s3_handler - INFO - Uploaded /data/MEMORANDUM OF UNDERSTANDING BETWEEN LACPD BERKELEY (CLJC edits).pdf to s3://cpda-documents/memorandum/2025/05/01/MEMORANDUM OF UNDERSTANDING BETWEEN LACPD BERKELEY (CLJC edits)_d17c5fda.pdf
app-1  | 2025-05-01 00:22:10,417 - document_processor - INFO - Uploaded to S3: s3://cpda-documents/memorandum/2025/05/01/MEMORANDUM OF UNDERSTANDING BETWEEN LACPD BERKELEY (CLJC edits)_d17c5fda.pdf
Processing files:  93%|█████████▎| 14/15 [00:14<00:00,  1.29it/s]2025-05-01 00:22:10,938 - document_classifier - INFO - LLM classification/extraction successful: Type 'Motion', Tags: ['Motion to Dismiss', 'Juvenile Adjudication', 'Strike Offenses', 'Penal Code Violations', 'Aiding and Abetting'], Timestamp: 'None'
app-1  | 2025-05-01 00:22:10,938 - document_processor - INFO - Document classified as: Motion, Timestamp: None, Legal Tags: ['Motion to Dismiss', 'Juvenile Adjudication', 'Strike Offenses', 'Penal Code Violations', 'Aiding and Abetting']
app-1  | 2025-05-01 00:22:11,502 - s3_handler - INFO - Uploaded /data/707(B).pdf to s3://cpda-documents/motion/2025/05/01/707(B)_fb12a4a6.pdf
app-1  | 2025-05-01 00:22:11,503 - document_processor - INFO - Uploaded to S3: s3://cpda-documents/motion/2025/05/01/707(B)_fb12a4a6.pdf
Processing files: 100%|██████████| 15/15 [00:15<00:00,  1.01s/it]
app-1  | 2025-05-01 00:22:12,847 - elastic_transport.transport - INFO - PUT https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/_bulk [status:200 duration:1.339s]
app-1  | 2025-05-01 00:22:12,847 - elasticsearch_handler - INFO - Bulk indexing complete. Success: 10, Errors: 0
app-1  | 2025-05-01 00:22:12,847 - document_processor - INFO - Processing complete. Stats: {'total': 15, 'processed': 10, 'indexed': 10, 'uploaded': 10, 'skipped': 3, 'failed': 4, 'unsupported_format': 2, 'too_large': 0}
app-1  |
app-1  | Processing Summary:
app-1  | Total files found: 15
app-1  | Files processed: 10
app-1  | Documents indexed: 10
app-1  | Files uploaded to S3: 10
app-1  | Files skipped: 3
app-1  | Failed files: 4
app-1 exited with code 0