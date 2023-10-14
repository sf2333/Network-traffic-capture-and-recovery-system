#include "AhoCorasick.h"


AcNode::AcNode() {
	fail = nullptr;
	memset(next, '\0', sizeof(next));
	count = 0;
	str = '\0';
}

//建立时插入字符串，构建树
void AhoCoraSick::Insert(byte* str, AcNode* root,  int identity) {
	auto p = root;

	for (auto i = 0; str[i] != '\0'; i++) {
		const auto index = str[i];

		if ( p->next[index]== nullptr ) {
			const auto q = new AcNode();
			/*q->str = index;*/
			p->next[index] = q;
			p = q;
		}
		else {
			p = p->next[index];
		}
	}
	p->count = identity;
}

//构建fail
void AhoCoraSick::BuildAcFail(AcNode* root) {
	if (root == nullptr) {
		return;
	}

	root->fail = root;

	queue<AcNode*> ac_queue;


	for (auto& i : root->next) {
		if (i != nullptr) {
			i->fail = root;
			ac_queue.push(i);
		}

	}


	while (!ac_queue.empty()) {
		auto p = ac_queue.front();
		ac_queue.pop();

		for (auto i = 0; i < 256; i++) {

			if (p->next[i] != nullptr) {
				if (p->fail->next[i] != nullptr) {
					p->next[i]->fail = p->fail->next[i];
				}
				else {
					p->next[i]->fail = p->fail;
				}
				ac_queue.push(p->next[i]);
			}

		}

	}
}


int AhoCoraSick::SearchKeyword(const byte* buffer, const int length, AcNode* root) {


	auto p = root;
	auto number = 0;


	for (auto i = 0; i < length; i++) {
		const auto str = buffer[i];

		while ((p->next[str]== nullptr) && (p != root)) {
			p = p->fail;
		}

		if ( p->next[str] != nullptr ) {
			p = p->next[str];
			if (p->count > 0) {
				number++;
				cout << "位置为：" << i << endl;
			}
		}
	}

	return number;
}