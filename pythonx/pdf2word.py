from pdf2docx import Converter

pdf_file = './assets/1320合格证.pdf'
docx_file = '1320合格证.docx'

pdf_file = './assets/新建 DOCX 文档.pdf'
docx_file = '1.docx'

# convert pdf to docx
cv = Converter(pdf_file)
cv.convert(docx_file)  # all pages by default
cv.close()
