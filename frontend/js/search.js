// Search functionality with debounce
const Search = {
    timeout: null,
    delay: 300, // milliseconds

    init() {
        const $searchInput = $('#searchInput');
        
        $searchInput.on('input', () => {
            clearTimeout(this.timeout);
            this.timeout = setTimeout(() => {
                const query = $searchInput.val().trim();
                CustomerTable.currentPage = 1; // Reset to first page on new search
                CustomerTable.render(query);
            }, this.delay);
        });

        $('#clearSearchBtn').on('click', () => {
            $searchInput.val('').focus();
            CustomerTable.currentPage = 1;
            CustomerTable.render('');
        });
    }
};
