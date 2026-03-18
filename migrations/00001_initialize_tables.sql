-- +goose Up
-- +goose StatementBegin
create table public.units
(
    id          uuid default uuidv7() not null
        constraint units_id_pk
            primary key,
    title       varchar(128)          not null,
    description varchar(1024)         not null
);

create table public.shipments
(
    id               uuid      default uuidv7() not null
        constraint shipments_id_pk
            primary key,
    reference_number varchar(64)                not null,
    unit_id          uuid                       not null
        constraint shipments_units_id_fk
            references public.units,
    origin           varchar(256)               not null,
    destination      varchar(256)               not null,
    driver_name      varchar(128)               not null,
    shipment_cost    bigint                     not null,
    driver_revenue   bigint                     not null,
    created_at       timestamp default now()    not null
);

create table public.events
(
    id          uuid      default uuidv7() not null
        constraint events_id_pk
            primary key,
    shipment_id uuid                       not null
        constraint events_shipments_id_fk
            references public.shipments,
    status      integer                    not null,
    details     varchar(256)               not null,
    occurred_at timestamp default now()    not null
);

INSERT INTO units (id, title, description)
VALUES ('019cff62-d5f7-7467-aa5f-dbc4497da6b8', 'Ergonomic Office Chair', 'High-back mesh chair with adjustable lumbar support and 4D armrests.'),
       ('019cff62-d5f9-746f-a5a2-bcc9ff0fc39d', 'Portable Power Bank 20k', '20,000mAh external battery with 65W PD fast charging and dual USB-C ports.'),
       ('019cff62-d5f9-74d9-866d-0dc2702c54ee', 'Wireless Noise-Canceling Headphones', 'Over-ear Bluetooth headphones with 40-hour battery life and active ANC.'),
       ('019cff62-d5f9-74e6-803e-108542f38459', 'Smart Coffee Maker', 'Wi-Fi enabled drip coffee machine with programmable brew strength and timer.'),
       ('019cff62-d5f9-74ef-be69-025d60153822', 'Stainless Steel Cookware Set', '10-piece professional-grade induction-ready pots and pans with glass lids.'),
       ('019cff62-d5f9-74f7-bd7b-aab361fbaefe', 'Electric Standing Desk', 'Motorized height-adjustable desk with memory presets and anti-collision sensor.'),
       ('019cff62-d5f9-74fe-993d-ac18f03648f8', 'HEPA Air Purifier', 'Compact 3-stage filtration system for rooms up to 300 sq. ft. with CADR 250.'),
       ('019cff62-d5f9-7507-b46f-ad806b95d32f', 'Resistance Band Set', '11-piece home workout kit including 5 stackable bands, handles, and door anchor.'),
       ('019cff62-d5f9-750e-a8b2-744652b8059e', 'Outdoor Camping Tent', '4-person waterproof dome tent with rainfly and easy-setup fiberglass poles.'),
       ('019cff62-d5f9-7516-8dbc-2a03941959dc', 'Mechanical Gaming Keyboard', 'RGB backlit wired keyboard with hot-swappable switches and PBT keycaps.');

INSERT INTO public.shipments (id, reference_number, unit_id, origin, destination, driver_name, shipment_cost, driver_revenue)
VALUES ('019cff8f-41e6-79f8-b6ef-728d33f1990a', 'DH5Q3E0J5X2E', '019cff62-d5f7-7467-aa5f-dbc4497da6b8', 'Chicago, IL', 'New York, NY', 'John Doe', 120000, 95000),
       ('019cff8f-41e8-71d0-9431-a472bf32f48e', 'DH5QAN25KN69', '019cff62-d5f9-746f-a5a2-bcc9ff0fc39d', 'Los Angeles, CA', 'Phoenix, AZ', 'Sarah Smith', 85000, 68000),
       ('019cff8f-41e8-7271-a375-1c6a7a84cdf4', 'DH5OOHTPQ4TR', '019cff62-d5f9-74d9-866d-0dc2702c54ee', 'Seattle, WA', 'Portland, OR', 'Mike Johnson', 45000, 36000),
       ('019cff8f-41e8-728b-9b77-34a660c66186', 'DH5RHJKETRUN', '019cff62-d5f9-74e6-803e-108542f38459', 'Austin, TX', 'Dallas, TX', 'Elena Rodriguez', 50000, 40000),
       ('019cff8f-41e8-7297-b7b6-74ce7f1a8cf8', 'DH5RX0RFHDAQ', '019cff62-d5f9-74ef-be69-025d60153822', 'Atlanta, GA', 'Miami, FL', 'David Chen', 110000, 88000),
       ('019cff8f-41e8-72a3-891a-3d59c200df16', 'DH5RRBX6KF8G', '019cff62-d5f9-74f7-bd7b-aab361fbaefe', 'Denver, CO', 'Salt Lake City, UT', 'Chris Evans', 95000, 76000),
       ('019cff8f-41e8-72ad-819d-4339dfd948d3', 'DH5PKAG2AUYL', '019cff62-d5f9-74fe-993d-ac18f03648f8', 'Boston, MA', 'Philadelphia, PA', 'Lisa Wong', 60000, 48000),
       ('019cff8f-41e8-72b9-a8b3-50a5fee94c6a', 'DH5RT68HWEMC', '019cff62-d5f9-7507-b46f-ad806b95d32f', 'Nashville, TN', 'Charlotte, NC', 'Kevin Hart', 70000, 56000),
       ('019cff8f-41e8-72c3-8fee-979d16a47ba3', 'DH5PGKRZ4BVX', '019cff62-d5f9-750e-a8b2-744652b8059e', 'San Francisco, CA', 'Las Vegas, NV', 'Maria Garcia', 130000, 104000),
       ('019cff8f-41e8-72cd-b1e4-47b268f143d0', 'DH5PIULEAYDU', '019cff62-d5f9-7516-8dbc-2a03941959dc', 'Detroit, MI', 'Columbus, OH', 'Robert Taylor', 40000, 32000);

INSERT INTO public.events (shipment_id, status, details, occurred_at)
VALUES ('019cff8f-41e6-79f8-b6ef-728d33f1990a', 0, 'Shipment created.', now()),
       ('019cff8f-41e6-79f8-b6ef-728d33f1990a', 1, 'Driver assigned, awaiting pickup in Chicago.', now() + INTERVAL '1 minute'),
       ('019cff8f-41e6-79f8-b6ef-728d33f1990a', 2, 'Picked up by John Doe.', now() + INTERVAL '2 minutes'),
       ('019cff8f-41e6-79f8-b6ef-728d33f1990a', 3, 'In transit to New York.', now() + INTERVAL '3 minutes'),
       ('019cff8f-41e6-79f8-b6ef-728d33f1990a', 6, 'Successfully delivered to destination.', now() + INTERVAL '4 minutes'),

       ('019cff8f-41e8-71d0-9431-a472bf32f48e', 0, 'Shipment created.', now() + INTERVAL '5 minutes'),
       ('019cff8f-41e8-71d0-9431-a472bf32f48e', 1, 'Awaiting driver Sarah Smith.', now() + INTERVAL '6 minutes'),
       ('019cff8f-41e8-71d0-9431-a472bf32f48e', 2, 'Picked up in Los Angeles.', now() + INTERVAL '7 minutes'),
       ('019cff8f-41e8-71d0-9431-a472bf32f48e', 3, 'In transit.', now() + INTERVAL '8 minutes'),
       ('019cff8f-41e8-71d0-9431-a472bf32f48e', 4, 'Delayed due to heavy traffic on I-10.', now() + INTERVAL '9 minutes'),

       ('019cff8f-41e8-7271-a375-1c6a7a84cdf4', 0, 'Shipment created.', now() + INTERVAL '10 minutes'),
       ('019cff8f-41e8-7271-a375-1c6a7a84cdf4', 1, 'Awaiting driver Mike Johnson.', now() + INTERVAL '11 minutes'),
       ('019cff8f-41e8-7271-a375-1c6a7a84cdf4', 7, 'Cancelled by customer before pickup.', now() + INTERVAL '12 minutes'),

       ('019cff8f-41e8-728b-9b77-34a660c66186', 0, 'Shipment created.', now() + INTERVAL '13 minutes'),
       ('019cff8f-41e8-728b-9b77-34a660c66186', 1, 'Awaiting driver Elena Rodriguez.', now() + INTERVAL '14 minutes'),
       ('019cff8f-41e8-728b-9b77-34a660c66186', 2, 'Picked up in Austin.', now() + INTERVAL '15 minutes'),
       ('019cff8f-41e8-728b-9b77-34a660c66186', 3, 'In transit.', now() + INTERVAL '16 minutes'),
       ('019cff8f-41e8-728b-9b77-34a660c66186', 5, 'Arrived at Dallas transfer hub.', now() + INTERVAL '17 minutes'),

       ('019cff8f-41e8-7297-b7b6-74ce7f1a8cf8', 0, 'Shipment created.', now() + INTERVAL '18 minutes'),
       ('019cff8f-41e8-7297-b7b6-74ce7f1a8cf8', 1, 'Awaiting driver David Chen.', now() + INTERVAL '19 minutes'),
       ('019cff8f-41e8-7297-b7b6-74ce7f1a8cf8', 2, 'Picked up in Atlanta.', now() + INTERVAL '20 minutes'),

       ('019cff8f-41e8-72a3-891a-3d59c200df16', 0, 'Shipment created.', now() + INTERVAL '21 minutes'),
       ('019cff8f-41e8-72a3-891a-3d59c200df16', 1, 'Awaiting driver Chris Evans.', now() + INTERVAL '22 minutes'),
       ('019cff8f-41e8-72a3-891a-3d59c200df16', 2, 'Picked up in Denver.', now() + INTERVAL '23 minutes'),
       ('019cff8f-41e8-72a3-891a-3d59c200df16', 3, 'In transit to Salt Lake City.', now() + INTERVAL '24 minutes'),

       ('019cff8f-41e8-72ad-819d-4339dfd948d3', 0, 'Shipment created.', now() + INTERVAL '25 minutes'),
       ('019cff8f-41e8-72ad-819d-4339dfd948d3', 1, 'Confirmed, awaiting Lisa Wong for pickup.', now() + INTERVAL '26 minutes'),

       ('019cff8f-41e8-72b9-a8b3-50a5fee94c6a', 0, 'Initial request received from Nashville.', now() + INTERVAL '27 minutes'),

       ('019cff8f-41e8-72c3-8fee-979d16a47ba3', 0, 'Shipment created.', now() + INTERVAL '28 minutes'),
       ('019cff8f-41e8-72c3-8fee-979d16a47ba3', 1, 'Awaiting driver Maria Garcia.', now() + INTERVAL '29 minutes'),
       ('019cff8f-41e8-72c3-8fee-979d16a47ba3', 2, 'Picked up.', now() + INTERVAL '30 minutes'),
       ('019cff8f-41e8-72c3-8fee-979d16a47ba3', 3, 'In transit.', now() + INTERVAL '31 minutes'),
       ('019cff8f-41e8-72c3-8fee-979d16a47ba3', 4, 'Weather delay in Sierra Nevada.', now() + INTERVAL '32 minutes'),
       ('019cff8f-41e8-72c3-8fee-979d16a47ba3', 3, 'Resumed transit to Las Vegas.', now() + INTERVAL '33 minutes'),

       ('019cff8f-41e8-72cd-b1e4-47b268f143d0', 0, 'Shipment created.', now() + INTERVAL '34 minutes'),
       ('019cff8f-41e8-72cd-b1e4-47b268f143d0', 1, 'Awaiting driver Robert Taylor.', now() + INTERVAL '35 minutes'),
       ('019cff8f-41e8-72cd-b1e4-47b268f143d0', 2, 'Picked up in Detroit.', now() + INTERVAL '36 minutes'),
       ('019cff8f-41e8-72cd-b1e4-47b268f143d0', 3, 'In transit.', now() + INTERVAL '37 minutes'),
       ('019cff8f-41e8-72cd-b1e4-47b268f143d0', 6, 'Delivered to Columbus facility.', now() + INTERVAL '38 minutes');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table public.events;
drop table public.shipments;
drop table public.units;
-- +goose StatementEnd