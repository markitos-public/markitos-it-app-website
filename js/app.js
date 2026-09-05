const faqSearch = document.querySelector('#faq-search-input');
const faqItems = document.querySelectorAll('.faq-list details');
const faqEmptyState = document.querySelector('.faq-empty');

if (faqItems.length > 0) {
    faqItems.forEach((faq) => {
        faq.addEventListener('toggle', () => {
            if (!faq.open) {
                return;
            }

            faqItems.forEach((otherFaq) => {
                if (otherFaq !== faq) {
                    otherFaq.open = false;
                }
            });
        });
    });
}

if (faqSearch) {
    faqSearch.addEventListener('input', () => {
        const query = faqSearch.value.trim().toLowerCase();
        let visibleFaqs = 0;

        faqItems.forEach((faq) => {
            const searchableText = [faq.dataset.title, faq.dataset.content, faq.dataset.tags].join(' ').toLowerCase();
            const matches = !query || searchableText.includes(query);

            faq.hidden = !matches;
            if (matches) {
                visibleFaqs += 1;
            } else {
                faq.open = false;
            }
        });

        faqEmptyState.hidden = visibleFaqs > 0;
    });
}
