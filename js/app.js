const menuToggle = document.querySelector('.menu-toggle');
const siteNav = document.querySelector('.site-nav');
const filters = document.querySelectorAll('.topic-filter');
const cards = document.querySelectorAll('.article-card, .featured-article');
const searchInput = document.querySelector('#article-search');
const emptyState = document.querySelector('.empty-state');
const newsletterForm = document.querySelector('#newsletter-form');
const formMessage = document.querySelector('#form-message');

menuToggle.addEventListener('click', () => {
    const isOpen = siteNav.classList.toggle('open');
    menuToggle.setAttribute('aria-expanded', String(isOpen));
});

function updateArticles() {
    const activeFilter = document.querySelector('.topic-filter.active').dataset.filter;
    const searchTerm = searchInput.value.trim().toLowerCase();
    let visibleCards = 0;

    cards.forEach((card) => {
        const matchesFilter = activeFilter === 'all' || card.dataset.category === activeFilter;
        const matchesSearch = !searchTerm || card.dataset.title?.toLowerCase().includes(searchTerm) || card.textContent.toLowerCase().includes(searchTerm);
        const isVisible = matchesFilter && matchesSearch;
        card.hidden = !isVisible;
        if (isVisible) visibleCards += 1;
    });

    emptyState.hidden = visibleCards > 0;
}

filters.forEach((filter) => {
    filter.addEventListener('click', () => {
        filters.forEach((item) => item.classList.remove('active'));
        filter.classList.add('active');
        updateArticles();
    });
});

searchInput.addEventListener('input', updateArticles);

newsletterForm.addEventListener('submit', (event) => {
    event.preventDefault();
    formMessage.textContent = 'Listo. Revisa tu email para confirmar la suscripcion.';
    newsletterForm.reset();
});
