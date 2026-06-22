import re

def normalize_token(token: str) -> str:
    token = token.strip().lower()
    token = token.strip('"').strip("'")
    token = re.sub(r'[^a-z0-9]+', '', token)
    return token


def process_stopwords_to_file(input_path: str, output_path: str, dedupe: bool = True):
    seen = set()

    with open(input_path, "r", encoding="utf-8") as fin, \
         open(output_path, "w", encoding="utf-8") as fout:

        for line in fin:
            # split by comma or newline style blobs
            parts = re.split(r'[,\n]+', line)

            for p in parts:
                token = normalize_token(p)

                if not token:
                    continue

                # optional dedupe
                if dedupe:
                    if token in seen:
                        continue
                    seen.add(token)

                fout.write(token + "\n")
    
raw = '''

'''

# print(to_vertical_stopword_list(raw))

process_stopwords_to_file("in.txt", "out.txt", dedupe=False)