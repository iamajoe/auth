alter table {{ index .Options "Namespace" }}.mfa_factors drop constraint if exists mfa_factors_phone_key;
do $$
begin
    if exists (
         select 1
         from pg_indexes
         where indexname = 'unique_verified_phone_factor'
         and schemaname = '{{ index .Options "Namespace" }}'
    ) then
        -- Ensure the target index 'unique_phone_factor_per_user' doesn't already exist
        if not exists (
            select 1
            from pg_indexes
            where indexname = 'unique_phone_factor_per_user'
            and schemaname = '{{ index .Options "Namespace" }}'
        ) then
            -- Rename the index if the target name doesn't already exist
            execute 'alter index {{ index .Options "Namespace" }}.unique_verified_phone_factor rename to unique_phone_factor_per_user';
        end if;
    end if;
end $$;
