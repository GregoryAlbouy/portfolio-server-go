INSERT INTO project (
    name,
    slug,
    description,
    tags,
    image,
    repo,
    demo,
    is_hidden,
    added_on,
    edited_on
) VALUES (
    :name,
    :slug,
    :description,
    :tags,
    :image,
    :repo,
    :demo,
    :is_hidden,
    :added_on,
    :edited_on
);