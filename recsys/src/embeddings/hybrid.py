class HybridEmbedding:
    def __init__(self, dense, sparse):
        self.dense = dense
        self.sparse = sparse

    def encode_document(self, text):

        return {
            "dense": self.dense.encode_document(text),
            "sparse": self.sparse.encode_document(text),
        }
