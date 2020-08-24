UPDATE project
SET
    name=:name,
    slug=:slug,
    description=:description,
    image=:image,
    tags=:tags,
    repo=:repo,
    demo=:demo,
    is_hidden=:is_hidden,
    edited_on=:edited_on
WHERE id=:id