import React from 'react';
import SearchAndFilterAccommodations
    from "../../components/search-and-filter-accommodations/search-and-filter-accommodations";

function SearchAndFilterAccommodationsPage(props) {

    return (
        <div>
            <h1>Search And Filter Accommodations</h1>
            <SearchAndFilterAccommodations canBuy={props.canBuy}></SearchAndFilterAccommodations>
        </div>
    );
}

export default SearchAndFilterAccommodationsPage;