CREATE TABLE banner (
    id serial primary key,
    document text not null,
    created_at timestamp not null,
    updated_at timestamp,
    feature_id int not null,
    is_active boolean,
    CONSTRAINT unique_banner_doc UNIQUE (document)
);

CREATE TABLE banner_feature_tag (
    tag_id int,
    banner_id int,
    feature_id int,
    CONSTRAINT fk_banner FOREIGN KEY (banner_id) REFERENCES banner(id),
    CONSTRAINT unique_feature_tag UNIQUE (feature_id, tag_id)
    CONSTRAINT unique_feature_tag UNIQUE (banner_id, tag_id)
);

