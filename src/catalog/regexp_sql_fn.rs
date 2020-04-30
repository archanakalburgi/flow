use regex::Regex;
use rusqlite::{functions::FunctionFlags, Connection, Error, Result};

pub fn install(db: &Connection) -> Result<()> {
    let flags = FunctionFlags::SQLITE_UTF8 | FunctionFlags::SQLITE_DETERMINISTIC;
    db.create_scalar_function("regexp", 2, flags, move |ctx| {
        assert_eq!(ctx.len(), 2, "called with unexpected number of arguments");

        let saved_re: Option<&Regex> = ctx.get_aux(0)?;
        let new_re = match saved_re {
            None => {
                let s = ctx.get::<String>(0)?;
                match Regex::new(&s) {
                    Ok(r) => Some(r),
                    Err(err) => return Err(Error::UserFunctionError(Box::new(err))),
                }
            }
            Some(_) => None,
        };

        let is_match = {
            let re = saved_re.unwrap_or_else(|| new_re.as_ref().unwrap());

            let text = ctx
                .get_raw(1)
                .as_str()
                .map_err(|e| Error::UserFunctionError(e.into()))?;

            re.is_match(text)
        };

        if let Some(re) = new_re {
            ctx.set_aux(0, re);
        }

        Ok(is_match)
    })
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_function() -> Result<()> {
        let db = Connection::open_in_memory()?;
        install(&db)?;

        let is_match: bool = db.query_row(
            "SELECT 'aaaaeeeiii' REGEXP '^[aeiou]*$';",
            rusqlite::NO_PARAMS,
            |row| row.get(0),
        )?;

        assert!(is_match);
        Ok(())
    }

    #[test]
    fn test_install_twice() -> Result<()> {
        // Can be installed multiple times without issue.
        let db = Connection::open_in_memory()?;
        install(&db)?;
        install(&db)?;
        Ok(())
    }
}